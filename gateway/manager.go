package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"go-im/common/conf/serviceConf"
	"go-im/common/discovery"
	"go-im/common/log"
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"go-im/common/util"
	"go-im/gateway/epoll"
	"go-im/gateway/rpc/client"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type iManager interface {
}

type Manager struct {
	deviceId   int32                      //设备编号
	addr       string                     //服务地址
	lis        *net.TCPListener           //tcp监听器
	table      sync.Map                   //连接和id的映射
	workPool   *ants.Pool                 //协程池
	hasConnNum int32                      //现有连接数
	connChan   chan *Connection           //accept conn事件后，通知epoll
	recv       chan *MessageEvent         // 处理上游消息
	send       chan *MessageEvent         //处理下游消息
	register   *discovery.ServiceRegister //服务注册
	done       chan struct{}              //manager关闭时，done关闭
}

var m *Manager

func initManager() error {
	m = &Manager{
		deviceId: serviceConf.GetGateWayDeviceId(),
		addr:     serviceConf.GetGateWayAddr(),
		connChan: make(chan *Connection, 3),
		recv:     make(chan *MessageEvent, 10),
		send:     make(chan *MessageEvent, 10),
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", m.addr)
	if err != nil {
		return err
	}
	lis, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	m.lis = lis
	if m.workPool, err = ants.NewPool(serviceConf.GetGateWayWorkPoolNum()); err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	m.accept()
	m.initEpoll()
	//处理上游消息
	go m.handleUpstreamMessage()
	// 处理下游消息
	go m.handleDownstreamMessage()
	log.Info("======================= IM Gateway ===================== ")
	//注册服务到etcd
	m.registerService()
	return nil
}

func (m *Manager) accept() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, err := m.lis.AcceptTCP()
				if err != nil {
					log.Error("failed accept: ", err)
					return
				}
				c := codec.NewGobCodec(conn)
				account, err := m.auth(c)
				if err != nil {
					log.Info(err)
					_ = c.Close()
					continue
				}
				connN := NewConnection(getSocketFd(conn), account, c, m)
				m.storeConn(connN)
				m.connChan <- connN
			}
		}()
	}
}

// 鉴权 + 限流
func (m *Manager) auth(codec codec.Codec) (string, error) {
	//读取鉴权信息
	reqH := &tcp.FixedHeader{}
	reqBody := &tcp.AuthMB{}
	respH := &tcp.FixedHeader{MsgId: -1, PreMsgId: -1, MessageType: tcp.AuthResponseMessage}
	respBody := &tcp.AuthResponseMB{}
	if err := codec.ReadFixedHeader(reqH); err != nil {
		respBody.Status = -1
		respBody.ErrMsg = err.Error()
		_ = codec.Write(respH, respBody)
		return "", err
	}
	if err := codec.ReadBody(reqBody); err != nil {
		respBody.Status = -1
		respBody.ErrMsg = err.Error()
		_ = codec.Write(respH, respBody)
		return "", err
	}
	//限流判断
	if !m.CheckAndAddTcpNum() {
		err := errors.New("conn num has reached max,confuse connection")
		respBody.Status = 2
		respBody.ErrMsg = err.Error()
		_ = codec.Write(respH, respBody)
		return "", err
	}
	//调用限流rpc接口
	ok, err := client.Auth(reqH.From, reqBody.Token)
	//回复响应
	if !ok {
		//鉴权失败
		respBody.Status = 0
		respBody.ErrMsg = err.Error()
		_ = codec.Write(respH, respBody)
		return "", err
	}
	respBody.Status = 1
	return respH.From, codec.Write(respH, respBody)
}

func (m *Manager) initEpoll() {
	for i := 0; i < serviceConf.GetGateWayReactorNums(); i++ {
		go m.StartEpoll()
	}
}

func (m *Manager) StartEpoll() {
	e, err := epoll.CreateEpoll()
	if err != nil {
		log.Info(err)
	}
	//开启协程监听accept事件，并将连接注册到epoll
	go func() {
		for {
			select {
			case <-m.done:
				log.Infof("epoll %d closed", e.Epollfd)
				_ = e.Close()
				return
			case conn := <-m.connChan:
				log.Infof("epoll %d 添加连接 %d", e.Epollfd, conn.Fd)
				conn.EpollFd = e.Epollfd
				err := e.AddEpollTask(int32(conn.Fd))
				if err != nil {
					log.Info(err)
				}
			}
		}
	}()
	// 处理触发事件
	for {
		epollEvent, n, err := e.EventTigger()
		if err != nil && err != syscall.EINTR {
			// EINTR 表示遇到中断，这里原因是epoll_wait超时
			// TODO 处理其他错误
			log.Info(err)
			continue
		}
		for i := 0; i < n; i++ {
			conn, err := m.load(int(epollEvent[i].Fd))
			if err != nil {
				log.Error(err)
				continue
			}
			m.ReadMessage(conn, e)
		}
	}
}

func (m *Manager) ReadMessage(conn *Connection, e *epoll.Epoller) {
	h := &tcp.FixedHeader{}
	if err := conn.Codec.ReadFixedHeader(h); err != nil {
		// 连接关闭
		m.Remove(conn.Fd)
		if err := e.DelEpollTask(int32(conn.Fd)); err != nil {
			log.Error(e.Epollfd, "failed to del", conn.Fd, err)
		}
		_ = conn.Close()
		return
	}
	body := tcp.GetMessageBody(h.MessageType)
	if err := conn.Codec.ReadBody(body); err != nil {
		log.Error(err)
		return
	}
	m.recv <- NewMessageEvent(conn.Fd, h, body)
}

// CheckAndAddTcpNum 检查并添加连接数
func (m *Manager) CheckAndAddTcpNum() bool {
	//已经超过最大连接数
	if atomic.LoadInt32(&m.hasConnNum) >= serviceConf.GetGateWayMaxConnsNum() {
		return false
	}
	return true
}

func (m *Manager) storeConn(conn *Connection) {
	atomic.AddInt32(&m.hasConnNum, 1)
	m.table.Store(conn.Fd, conn)

}

// Remove 删除连接在m上的映射
func (m *Manager) Remove(id int) {
	log.Infof("Remove a connection %d", id)
	atomic.AddInt32(&m.hasConnNum, -1)
	m.table.Delete(id)
}

func (m *Manager) load(id int) (*Connection, error) {
	if value, ok := m.table.Load(id); ok {
		return value.(*Connection), nil
	}
	return nil, errors.New(fmt.Sprintf("not found user %d", id))
}

func getSocketFd(conn *net.TCPConn) int {
	connV := reflect.ValueOf(*conn)
	fdV := reflect.Indirect(connV.FieldByName("fd"))
	return int(fdV.FieldByName("pfd").FieldByName("Sysfd").Int())
}

func (m *Manager) UpstreamMessageTask(messE *MessageEvent) func() {
	return func() {
		switch messE.Header.MessageType {
		case tcp.PrivateChatMessage:
			m.send <- messE
		case tcp.HeartBeatMessage:
			m.handleHeartBeatMB(messE)
		}
	}
}

func (m *Manager) handleUpstreamMessage() {
	for messE := range m.recv {
		for {
			err := m.workPool.Submit(m.UpstreamMessageTask(messE))
			if err != nil {
				if err == ants.ErrPoolOverload {
					continue
				}
				log.Error(err)
			}
			break
		}
	}
}

// 利用闭包解决进行协程池参数的传递
func (m *Manager) downstreamMessageTask(conn *Connection, messE *MessageEvent) func() {
	return func() {
		err := conn.Codec.Write(messE.Header, messE.Body)
		if err != nil {
			log.Error(err)
			_ = conn.Close()
			m.Remove(conn.Fd)
		}
	}
}

func (m *Manager) handleDownstreamMessage() {
	for messE := range m.send {
		conn, err := m.load(messE.connId)
		if err != nil {
			log.Errorf("user %d not exist on this gateway", messE.connId)
			continue
		}
		for {
			err = m.workPool.Submit(m.downstreamMessageTask(conn, messE))
			if err != nil {
				if err == ants.ErrPoolOverload {
					continue
				}
				log.Error(err)
			}
			break
		}
	}
}

func (m *Manager) registerService() {
	m.register = discovery.NewServerRegister(
		context.Background(),
		serviceConf.GetGatewayEndPoints(),
		serviceConf.GetGatewayDailTimeOut(),
		serviceConf.GetGatewayLeaseDDL(),
		fmt.Sprintf("im/gatewayServer/%d", m.deviceId),
		discovery.Transform(m.addr, m.getConnNums(), util.CPUPercent()))
	go m.regularUpdateService()
	go m.register.ListenKeepAliveChan()

}

//定期更新服务
func (m *Manager) regularUpdateService() {
	for {
		tick := time.NewTicker(time.Second)
		select {
		case <-m.done:
			tick.Stop()
			return
		case <-tick.C:
			m.register.UpdateService(map[string]interface{}{"connect_num": m.getConnNums(), "message_bytes": util.CPUPercent()})
		}
	}
}

func (m *Manager) getConnNums() float64 {
	return float64(atomic.LoadInt32(&m.hasConnNum))
}

func (m *Manager) Close() {
	m.workPool.Release()
	close(m.send)
	close(m.connChan)
	m.table.Range(func(key, value any) bool {
		conn := value.(*Connection)
		_ = conn.Close()
		return true
	})
}

// handleHeartBeatMB 处理心跳消息
func (m *Manager) handleHeartBeatMB(msg *MessageEvent) {
	conn, err := m.load(msg.connId)
	if err != nil {
		log.Info(err)
		return
	}
	conn.resetHeartBeatTimer()
}
