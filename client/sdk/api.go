package sdk

import (
	"encoding/json"
	"go-im/common/conf/serviceConf"
	"go-im/common/tcp"
	"io"
	"net/http"
)

type Chat struct {
	Nick           string
	UserID         string
	SessionID      string
	conn           *connect
	serverAddrList []string
}

type Message struct {
	Header *tcp.FixedHeader
	Body   interface{}
}

func NewChat(nick, userID, sessionID string) *Chat {
	c := &Chat{
		Nick:           nick,
		UserID:         userID,
		SessionID:      sessionID,
		serverAddrList: LoadBalanceIpList(),
	}
	c.conn = newConnet(c.serverAddrList)
	return c
}

func LoadBalanceIpList() []string {
	resp, err := http.Get(serviceConf.GetClientDiscoveryAddr())
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	ipListBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ipList := make([]string, 0)
	if err = json.Unmarshal(ipListBytes, &ipList); err != nil {
		panic(err)
	}
	return ipList
}

func (c *Chat) SendText(to string, t string) {
	h := &tcp.FixedHeader{
		Seq:         1,
		MessageType: tcp.PrivateChatMessage,
	}
	body := &tcp.PrivateChat{
		From:    c.UserID,
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
