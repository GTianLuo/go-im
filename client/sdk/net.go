package sdk

import (
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"go-im/log"
	"io"
	"net"
)

type connect struct {
	serverAddr         string
	codec              codec.Codec
	sendChan, recvChan chan *Message
}

func newConnet(serverAddr string) *connect {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	c := &connect{
		serverAddr: serverAddr,
		codec:      codec.NewGobCodec(conn),
		sendChan:   make(chan *Message),
		recvChan:   make(chan *Message),
	}
	go c.handleSendChan()
	go c.recvMessage()
	return c
}

func (c *connect) recvMessage() {
	for {
		h := &tcp.FixedHeader{}
		if err := c.codec.ReadFixedHeader(h); err != nil {
			// 连接关闭
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return
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
		c.recvChan <- &Message{Header: h, Body: body}
	}
}

func (c *connect) send(m *Message) {
	// 直接发送给接收方
	c.sendChan <- m
}

func (c *connect) getRecvChan() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
	// 目前没啥值得回收的
}

func (c *connect) handleSendChan() {
	for m := range c.sendChan {
		if err := c.codec.Write(m.Header, m.Body); err != nil {
			log.Fatal(err)
		}
	}
}
