package gateway

import (
	"go-im/common/conf/serviceConf"
	"go-im/common/log"
	"go-im/common/tcp/codec"
	"go-im/common/timingwheel"
	"go-im/gateway/epoll"
)

type Connection struct {
	Fd             int         // 文件描述符
	Account        string      //用户账号
	Codec          codec.Codec // 编解码器
	EpollFd        int         //epoll fd
	M              *Manager    // M 和EpollFd是为了conn主动删除连接
	heartBeatTimer *timingwheel.Timer
	msgTimer       *timingwheel.Timer
}

func NewConnection(fd int, account string, codec codec.Codec, m *Manager) *Connection {
	conn := &Connection{
		Fd:      fd,
		Codec:   codec,
		Account: account,
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

// CancelConn 关闭连接
func (conn *Connection) CancelConn() {
	conn.M.Remove(conn.Fd)
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
