package gateway

import (
	"go-im/common/conf/serviceConf"
	"go-im/common/dao"
	"go-im/common/log"
	"go-im/common/tcp/codec"
	"go-im/common/timingwheel"
	"go-im/gateway/epoll"
	"sync"
)

const (
	NoHandle = iota
	NeedAck
	NeedHandleAndAck
)

type Connection struct {
	mu             sync.Mutex
	Fd             int         // 文件描述符
	Account        string      //用户账号
	Codec          codec.Codec // 编解码器
	EpollFd        int         //epoll fd
	M              *Manager    // M 和EpollFd是为了conn主动删除连接
	MsgId          int64       //消息会话级别的唯一id
	heartBeatTimer *timingwheel.Timer
	msgTimer       *timingwheel.Timer
}

func NewConnection(fd int, account string, msgId int64, codec codec.Codec, m *Manager) *Connection {
	conn := &Connection{
		Fd:      fd,
		Codec:   codec,
		Account: account,
		MsgId:   msgId,
		M:       m,
	}
	conn.resetHeartBeatTimer()
	return conn
}

func (conn *Connection) Close() error {

	if conn.heartBeatTimer != nil {
		conn.heartBeatTimer.Stop()
	}
	if conn.msgTimer != nil {
		conn.msgTimer.Stop()
	}
	if err := conn.Codec.Close(); err != nil {
		return err
	}
	return nil
}

// SaveConnStatus 保存连接状态
func (conn *Connection) SaveConnStatus(deviceId int) error {
	return dao.NewGatewayStatus().SaveConnStatus(deviceId, conn.Account)
}

// DelConnStatus 删除连接状态
func (conn *Connection) DelConnStatus(deviceId int) error {
	return dao.NewGatewayStatus().DelConnStatus(deviceId, conn.Account)
}

// CancelConn 关闭连接
func (conn *Connection) CancelConn() {
	conn.M.removeConn(conn.Fd)
	if err := conn.DelConnStatus(int(conn.M.deviceId)); err != nil {
		log.Error(err)
	}
	epoller := epoll.GetEpoller(conn.EpollFd)
	if err := epoller.DelEpollTask(int32(conn.Fd)); err != nil {
		panic(err)
	}
	_ = conn.Close()
}

// resetHeartBeatTimer 重置心跳过期时间
func (conn *Connection) resetHeartBeatTimer() {

	if conn.heartBeatTimer != nil {
		conn.heartBeatTimer.Stop()
	}
	conn.heartBeatTimer = timingwheel.AfterFunc(serviceConf.GetGateWayHeartbeatTimeout(), func() {
		log.Infof("连接%d的心跳过期", conn.Fd)
		conn.CancelConn()
	})
}

// CheckAndAdd 通过MsgId检查消息的可靠性，不重不漏; 若MsgId无误，conn消耗一个msgId
func (conn *Connection) CheckAndAdd(msgId int64) int {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.MsgId == msgId {
		conn.MsgId++
		return NeedHandleAndAck
	} else if conn.MsgId > msgId {
		return NeedAck
	} else {
		return NoHandle
	}
}
