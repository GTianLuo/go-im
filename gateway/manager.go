package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"go-im/common/discovery"
	"go-im/common/tcp"
	"go-im/conf"
	"go-im/log"
	"io"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
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
	connChan   chan *connection           //accept conn事件后，通知epoll
	send       chan *MessageEvent         //处理下游消息
	register   *discovery.ServiceRegister //服务注册
	done       chan struct{}              //manager关闭时，done关闭
}

var m *Manager

func initManager() error {
	m = &Manager{
		deviceId: conf.GetGateWayDeviceId(),
		addr:     conf.GetGateWayAddr(),
		connChan: make(chan *connection, 3),
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
	if m.workPool, err = ants.NewPool(conf.GetGateWayWorkPoolNum()); err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	m.accept()
	m.initEpoll()
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
				// 限流
				if !m.CheckAndAddTcpNum() {
					log.Error("conn num has reached max,confuse connection")
					_ = conn.Close()
					continue
				}
				m.connChan <- NewConnection(getSocketFd(conn), conn)
			}
		}()
	}
}

func (m *Manager) initEpoll() {
	for i := 0; i < conf.GetGateWayReactorNums(); i++ {
		go m.StartEpoll()
	}
}

func (m *Manager) StartEpoll() {
	e, err := createEpoll()
	if err != nil {
		log.Info(err)
	}
	//开启协程监听accept事件，并将连接注册到epoll
	go func() {
		for {
			select {
			case <-m.done:
				log.Infof("epoll %d closed", e.epollfd)
				_ = e.close()
				return
			case conn := <-m.connChan:
				log.Infof("epoll %d 添加连接 %d", e.epollfd, conn.fd)
				err := e.addEpollTask(int32(conn.fd))
				if err != nil {
					log.Info(err)
				}
			}
		}
	}()
	// 处理触发事件
	for {
		epollEvent, n, err := e.eventTigger()
		if err != nil {
			log.Info(err)
			continue
		}
		for i := 0; i < n; i++ {
			connI, ok := m.table.Load(epollEvent[i].Fd)
			if !ok {
				continue
			}
			c := connI.(*connection)
			m.workPool.Submit(func() {
				h := &tcp.FixedHeader{}
				if err := c.codec.ReadFixedHeader(h); err != nil {
					// 连接关闭
					if err != nil {
						if err == io.EOF || err == io.ErrUnexpectedEOF {
							m.remove(c.fd)
							_ = e.delEpollTask(int32(e.epollfd))
						}
						// 读走坏数据
						_ = c.codec.ReadBody(nil)
						return
					}
				}
				body := tcp.GetMessageBody(h.MessageType)
				if err := c.codec.ReadBody(body); err != nil {
					log.Error(err)
					return
				}
				m.handleUpstreamMessage(NewMessageEvent(c.fd, h, body))
			})
		}
	}

}

// CheckAndAddTcpNum 检查并添加连接数
func (m *Manager) CheckAndAddTcpNum() bool {
	//已经超过最大连接数
	if atomic.LoadInt32(&m.hasConnNum) >= conf.GetGateWayMaxConnsNum() {
		return false
	}
	atomic.AddInt32(&m.hasConnNum, 1)
	return true
}

func (m *Manager) storeConn(conn *connection) {
	m.table.Store(conn.fd, conn)

}

func (m *Manager) remove(id int) {
	log.Infof("remove a connection %d", id)
	m.table.Delete(id)
}

func (m *Manager) load(id int) (*connection, error) {
	if value, ok := m.table.Load(id); ok {
		return value.(*connection), nil
	}
	return nil, errors.New(fmt.Sprintf("not found user %d", id))
}

func getSocketFd(conn *net.TCPConn) int {
	connV := reflect.ValueOf(*conn)
	fdV := reflect.Indirect(connV.FieldByName("fd"))
	return int(fdV.FieldByName("pfd").FieldByName("Sysfd").Int())
}

func (m *Manager) handleUpstreamMessage(event *MessageEvent) {
	m.send <- event
}

func (m *Manager) handleDownstreamMessage() {
	for cm := range m.send {
		conn, err := m.load(cm.UserID)
		if err != nil {
			log.Errorf("user %d not exist on this gateway", cm.UserID)
			continue
		}
		m.workPool.Submit(
			func() {
				err := conn.codec.Write(cm.Header, cm.Body)
				if err != nil {
					_ = conn.Close()
					m.remove(cm.UserID)
				}
			})
	}
}

func (m *Manager) registerService() {
	m.register = discovery.NewServerRegister(
		context.Background(),
		conf.GetGatewayEndPoints(),
		conf.GetGatewayDailTimeOut(),
		conf.GetGatewayLeaseDDL(),
		fmt.Sprintf("im/gatewayServer/%d", m.deviceId),
		discovery.Transform(m.addr, 0, 0))
}

func (m *Manager) Close() {
	m.workPool.Release()
	close(m.send)
	close(m.connChan)
	m.table.Range(func(key, value any) bool {
		conn := value.(*connection)
		_ = conn.Close()
		return true
	})
}
