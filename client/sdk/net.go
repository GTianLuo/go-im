package sdk

import (
	"errors"
	"go-im/common/log"
	"go-im/common/message"
	"go-im/common/tcp/codec"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type connect struct {
	serverAddr         string
	codec              codec.Codec
	sendChan, recvChan chan *message.Cmd
	closed             chan struct{}
}

func (chat *Chat) newConnet(serverAddrList []string) {
	c := &connect{
		sendChan: make(chan *message.Cmd, 10),
		recvChan: make(chan *message.Cmd, 10),
		closed:   make(chan struct{}, 1),
	}
	close(c.closed)
	chat.conn = c
	for {
		if err := c.login(serverAddrList, chat.account, chat.token, atomic.LoadInt64(&chat.MsgAckId)); err != nil {
			log.ClientError("failed to login:", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		return
	}
}

func (c *connect) login(serverAddrList []string, account string, token string, msgId int64) error {
	for _, serverAddr := range serverAddrList {
		// 建立连接
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			log.ClientError("failed dial to remote conn:", err.Error())
			continue
		}
		c.codec = codec.NewProtoCodec(conn)
		c.serverAddr = serverAddr
		//发送鉴权消息
		req := &message.AuthRequest{
			Token: token,
		}
		reqBytes, err := proto.Marshal(req)
		if err != nil {
			_ = conn.Close()
			panic(err)
		}
		authReqCmd := &message.Cmd{
			Type:    message.CmdType_AuthRequestCmd,
			MsgId:   msgId,
			From:    account,
			Payload: reqBytes,
		}
		if err = c.codec.WriteData(authReqCmd); err != nil {
			_ = conn.Close()
			return err
		}

		// 读取鉴权响应
		authRespCmd, err := c.codec.ReadData()
		if err != nil {
			_ = conn.Close()
			return err
		}
		authResp := &message.AuthResponse{}
		if err = proto.Unmarshal(authRespCmd.Payload, authResp); err != nil {
			panic(err)
		}
		if authResp.Status == 0 {
			c.close()
			panic(authResp.ErrMsg)
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
		cmd, err := c.codec.ReadData()
		if err != nil {
			// 连接关闭
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					//快速重连
					log.ClientInfo("Network disconnection")
					close(c.closed)
					c.reConnection()
					return
				}
				//非连接关闭的异常
				panic(err)
				return
			}
		}
		c.recvChan <- cmd
	}
}

// reConnection 尝试重新连接
func (c *connect) reConnection() {
	c.recvChan <- GetSystemMessage("连接断开，正在尝试重连", true)
}

func (c *connect) send(m *message.Cmd) {
	select {
	case <-c.closed:
		c.recvChan <- GetSystemMessage("无网络连接", false)
	default:
		c.sendChan <- m
	}
}

func (c *connect) getRecvChan() <-chan *message.Cmd {
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
			if err := c.codec.WriteData(m); err != nil {
				//log.Fatal(err)
				return
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
