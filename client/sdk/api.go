package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-im/common/conf/serviceConf"
	"go-im/common/result"
	"go-im/common/tcp"
	"io"
	"net/http"
	"strings"
	"time"
)

type Chat struct {
	nickName       string
	account        string
	token          string
	MsgId          int64
	conn           *connect
	serverAddrList []string
}

type Message struct {
	Header *tcp.FixedHeader
	Body   interface{}
}

func NewChat(account string, password string) (*Chat, error) {
	c := &Chat{}
	if err := c.LoadBalanceIpList(account, password); err != nil {
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
	fmt.Println(string(rBytes))
	r := &result.Result{}
	fmt.Println()
	if err = json.Unmarshal(rBytes, r); err != nil {
		panic(err)
	}
	// 登陆失败
	if r.Code != result.Ok {
		return errors.New(r.Msg)
	}
	user := r.Data.(map[string]interface{})
	for _, u := range user["ipList"].([]interface{}) {
		c.serverAddrList = append(c.serverAddrList, u.(string))
	}
	c.nickName = user["nickName"].(string)
	c.account = account
	c.token = user["token"].(string)
	return nil
}

func (c *Chat) SendText(to string, t string) {
	h := &tcp.FixedHeader{
		From:        c.account,
		PreMsgId:    c.MsgId,
		MsgId:       c.NextMsgId(),
		MessageType: tcp.PrivateChatMessage,
	}
	body := &tcp.PrivateChatMB{
		To:      to,
		Content: t,
	}
	c.conn.send(&Message{h, body})
}

//Close close chat
func (chat *Chat) Close() {
	chat.conn.close()
}

//Recv receive message
func (chat *Chat) Recv() <-chan *Message {
	return chat.conn.getRecvChan()
}

func (chat *Chat) NextMsgId() int64 {
	chat.MsgId = time.Now().UnixNano()
	return chat.MsgId
}
