package gateway

import (
	"go-im/common/tcp/codec"
	"net"
	"sync"
)

type connection struct {
	fd    int
	conn  *net.TCPConn
	codec codec.Codec
	send  sync.RWMutex
}

func NewConnection(fd int, conn *net.TCPConn) *connection {
	return &connection{
		fd:    fd,
		conn:  conn,
		codec: codec.NewGobCodec(conn),
	}
}

func (conn *connection) Close() error {
	return conn.conn.Close()
}
