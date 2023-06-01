package sdk

import (
	"errors"
	"go-im/common/log"
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"io"
	"net"
	"time"
)

type connect struct {
	serverAddr         string
	codec              codec.Codec
	sendChan, recvChan chan *Message
	closed             chan struct{}
}

func (chat *Chat) newConnet(serverAddrList []string) {
	c := &connect{
		sendChan: make(chan *Message, 10),
		recvChan: make(chan *Message, 10),
		closed:   make(chan struct{}, 1),
	}
	close(c.closed)
	chat.conn = c
	for {
		if err := c.login(serverAddrList, chat.account, chat.token); err != nil {
			//log.Info(err)
			time.Sleep(time.Second * 3)
			continue
		}
		return
	}
}

func (c *connect) login(serverAddrList []string, account string, token string) error {
	for _, serverAddr := range serverAddrList {
		// 建立连接
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			log.Error(err)
			continue
		}
		c.codec = codec.NewGobCodec(conn)
		c.serverAddr = serverAddr
		//发送鉴权消息
		header := &tcp.FixedHeader{
			MsgId:       -1,
			PreMsgId:    -1,
			From:        account,
			MessageType: tcp.AuthMessage,
		}
		body := &tcp.AuthMB{
			Token: token,
		}
		if err := c.codec.Write(header, body); err != nil {
			c.close()
			return err
		}
		// 读取鉴权响应
		respH := &tcp.FixedHeader{}
		respBody := &tcp.AuthResponseMB{}
		if err := c.codec.ReadFixedHeader(respH); err != nil {
			c.close()
			return err
		}
		if respH.MessageType != tcp.AuthResponseMessage {
			c.close()
			panic("wrong response type")
		}
		if err := c.codec.ReadBody(respBody); err != nil {
			c.close()
			return err
		}
		if respBody.Status == 0 {
			c.close()
			panic(respBody.ErrMsg)
		}
		// 登陆成功
		go c.handleSendChan()
		go c.recvMessage()
		go c.heartBeat(account)
		c.closed = make(chan struct{}, 1)
		return nil
	}
	return errors.New("No available connections")
}

func (c *connect) recvMessage() {
	for {
		h := &tcp.FixedHeader{}
		if err := c.codec.ReadFixedHeader(h); err != nil {
			// 连接关闭
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					//快速重连
					close(c.closed)
					c.reConnection()
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

// reConnection 尝试重新连接
func (c *connect) reConnection() {
	c.recvChan <- GetSystemMessage("连接断开，正在尝试重连", true)
}

func (c *connect) send(m *Message) {
	select {
	case <-c.closed:
		c.recvChan <- GetSystemMessage("无网络连接", false)
	default:
		c.sendChan <- m
	}
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
	for {
		select {
		case <-c.closed:
			return
		case m := <-c.sendChan:
			if err := c.codec.Write(m.Header, m.Body); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// heartBeat 心跳
func (c *connect) heartBeat(from string) {
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-c.closed:
			return
		case <-timer.C:
			c.sendChan <- GetHeartBeatMessage(from)
			timer.Reset(time.Second)
		}
	}
}
