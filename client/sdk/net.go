package sdk

import (
	"go-im/common/log"
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"io"
	"net"
)

type connect struct {
	serverAddr         string
	codec              codec.Codec
	sendChan, recvChan chan *Message
}

func (chat *Chat) newConnet(serverAddrList []string) {
	for i := 0; i < len(serverAddrList); i++ {
		//建立连接
		conn, err := net.Dial("tcp", serverAddrList[i])
		if err != nil {
			log.Error(err)
			continue
		}
		c := &connect{
			serverAddr: serverAddrList[i],
			codec:      codec.NewGobCodec(conn),
			sendChan:   make(chan *Message),
			recvChan:   make(chan *Message),
		}

		//发送鉴权消息
		header := &tcp.FixedHeader{
			MsgId:       -1,
			PreMsgId:    -1,
			From:        chat.account,
			MessageType: tcp.AuthMessage,
		}
		body := &tcp.AuthMB{
			Token: chat.token,
		}
		if err = c.codec.Write(header, body); err != nil {
			c.close()
			panic(err)
		}
		// 读取鉴权响应
		respH := &tcp.FixedHeader{}
		respBody := &tcp.AuthResponseMB{}
		if err = c.codec.ReadFixedHeader(respH); err != nil {
			c.close()
			panic(err)
		}
		if respH.MessageType != tcp.AuthResponseMessage {
			c.close()
			panic("wrong response type")
		}
		if err = c.codec.ReadBody(respBody); err != nil {
			c.close()
			panic(err)
		}
		if respBody.Status == 0 {
			c.close()
			panic(respBody.ErrMsg)
		}
		// 登陆成功
		chat.conn = c
		go c.handleSendChan()
		go c.recvMessage()
		return
	}
	panic("All connections are unavailable")
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
	_ = c.codec.Close()
	close(c.recvChan)
	close(c.sendChan)
}

func (c *connect) handleSendChan() {
	for m := range c.sendChan {
		if err := c.codec.Write(m.Header, m.Body); err != nil {
			log.Fatal(err)
		}
	}
}
