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
	"sync"
	"sync/atomic"
)

type iManager interface {
}

type Manager struct {
	deviceId   int32 //设备编号
	addr       string
	table      sync.Map           //连接和id的映射
	workPool   *ants.Pool         //协程池
	hasConnNum int32              //现有连接数
	recv       chan *MessageEvent //处理上游消息
	send       chan *MessageEvent //处理下游消息
	register   *discovery.ServiceRegister
}

var m *Manager

func initManager() error {
	m = &Manager{
		deviceId: conf.GetGateWayDeviceId(),
		addr:     conf.GetGateWayAddr(),
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
	if m.workPool, err = ants.NewPool(conf.GetGateWayWorkPoolNum()); err != nil {
		return errors.New("failed init Manager: " + err.Error())
	}
	go m.accept(lis)
	// 处理下游消息
	go m.handleDownstreamMessage()
	//处理上游消息
	go m.handleUpstreamMessage()
	log.Info("======================= IM Gateway ===================== ")
	//注册服务到etcd
	m.registerService()
	return nil
}

func (m *Manager) accept(lis *net.TCPListener) {
	for {
		conn, err := lis.AcceptTCP()
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
		go m.processConn(conn)
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

func (m *Manager) processConn(conn *net.TCPConn) {
	c := NewConnection(getSocketFd(conn), conn)
	defer func() { _ = c.Close() }()
	m.storeConn(c)
	log.Infof("create a new connection %d", c.fd)
	for {
		h := &tcp.FixedHeader{}
		if err := c.codec.ReadFixedHeader(h); err != nil {
			// 连接关闭
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					m.remove(c.fd)
				}
				// 读走坏数据
				_ = c.codec.ReadBody(nil)
				return
			}
		}
		body := tcp.GetMessageBody(h.MessageType)
		if err := c.codec.ReadBody(body); err != nil {
			log.Error(err)
			continue
		}
		m.recv <- NewMessageEvent(c.fd, h, body)
	}
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

func (m *Manager) handleUpstreamMessage() {
	for cm := range m.recv {
		m.send <- cm
	}
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
	close(m.recv)
	m.table.Range(func(key, value any) bool {
		conn := value.(*connection)
		_ = conn.Close()
		return true
	})
}
