package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"go-im/common/conf/serviceConf"
	"go-im/common/log"
	"go-im/common/message"
	"go-im/common/result"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Chat struct {
	mu             sync.Mutex
	nickName       string
	account        string
	token          string
	MsgId          int64
	MsgAckId       int64
	conn           *connect
	serverAddrList []string
	msgAckMap      map[int64]context.CancelFunc //收到ack后，cancel掉消息的重新
}

// NewChat 创建Chat对象，登陆账号获取token
func NewChat(account string, password string) (*Chat, error) {
	c := &Chat{MsgId: -1, MsgAckId: 0, msgAckMap: make(map[int64]context.CancelFunc)}
	if err := c.LoadBalanceIpList(account, password); err != nil {
		log.ClientError("login error:", err.Error())
		return c, err
	}
	c.newConnet(c.serverAddrList)
	return c, nil
}

// LoadBalanceIpList 登陆账号(获取登陆token)，并获取gateway网关
func (c *Chat) LoadBalanceIpList(account, password string) error {
	resp, err := http.Post(serviceConf.GetClientLoginAddr(), "application/json",
		strings.NewReader(`{"account":`+`"`+account+`",`+`"password":`+`"`+password+`"}`))
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	rBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r := &result.Result{}
	if err = json.Unmarshal(rBytes, r); err != nil {
		panic(err)
	}
	// 登陆失败
	if r.Code != result.Ok {
		return fmt.Errorf("%s:%s", r.Msg, r.Data)
	}
	user := r.Data.(map[string]interface{})
	for _, u := range user["ipList"].([]interface{}) {
		c.serverAddrList = append(c.serverAddrList, u.(string))
	}
	log.ClientInfof("success login； ipList: %v", c.serverAddrList)
	c.nickName = user["nickName"].(string)
	c.account = account
	c.token = user["token"].(string)
	return nil
}

// SendPrivateText 发送私聊消息
func (c *Chat) SendPrivateText(to string, t string) {
	msg := &message.PrivateMsg{
		Type: message.MsgType_TextMsg,
		To:   to,
		Data: []byte(t),
	}
	bytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	cmd := &message.Cmd{Type: message.CmdType_PrivateMsgCmd, MsgId: c.NextMsgId(), From: c.account, Payload: bytes}
	c.conn.send(cmd)
	go c.waitAck(cmd, t)
}

// 等待Ack，负责超时重传
func (c *Chat) waitAck(cmd *message.Cmd, content string) {
	ctx, cancel := context.WithCancel(context.Background())
	c.mu.Lock()
	c.msgAckMap[cmd.MsgId] = cancel
	c.mu.Unlock()
	for i := 1; i < serviceConf.GetClientMaxReSendNums(); i++ {
		select {
		case <-ctx.Done():
			//收到ack
			return
		case <-time.After(time.Millisecond * 500):
			// 超时重传
			log.ClientInfof("message.pb: %d 超时重传\n", cmd.MsgId)
			select {
			case <-c.conn.closed:
				continue
			default:
				c.conn.sendChan <- cmd
			}
		}
	}
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Second * 3):
		log.ClientInfof("message.pb: %d:%s 发送失败\n", cmd.MsgId, content)

		c.conn.recvChan <- GetSystemMessage("发送失败："+content, false)
	}
}

//Close close chat
func (chat *Chat) Close() {
	chat.conn.close()
}

//Recv receive message.pb
func (chat *Chat) Recv() <-chan *message.Cmd {
	return chat.conn.getRecvChan()
}

func (chat *Chat) NextMsgId() int64 {
	return atomic.AddInt64(&chat.MsgId, 1)
}

// ReConn 重连接
func (chat *Chat) ReConn() {
	for {
		if err := chat.conn.login(chat.serverAddrList, chat.account, chat.token, atomic.LoadInt64(&chat.MsgAckId)+1); err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		log.ClientInfo("reconnect success")
		return
	}
}

// HandleAck 处理ack消息
func (chat *Chat) HandleAck(msg *message.Cmd) {
	log.ClientInfof("收到message：%d 的ACK\n", msg.MsgId)
	chat.mu.Lock()
	if atomic.LoadInt64(&chat.MsgAckId) < msg.MsgId {
		atomic.StoreInt64(&chat.MsgAckId, msg.MsgId)
	}
	//取消超时重传
	chat.msgAckMap[msg.MsgId]()
	delete(chat.msgAckMap, msg.MsgId)
	chat.mu.Unlock()
}
