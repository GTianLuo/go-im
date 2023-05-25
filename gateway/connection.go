package gateway

import (
	"go-im/common/tcp/codec"
	"sync"
)

type connection struct {
	fd      int
	account string
	codec   codec.Codec
	send    sync.RWMutex
}

func NewConnection(fd int, codec codec.Codec) *connection {
	return &connection{
		fd:    fd,
		codec: codec,
	}
}

func (conn *connection) Close() error {
	return conn.codec.Close()
}
