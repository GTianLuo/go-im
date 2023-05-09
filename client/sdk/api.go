package sdk

import (
	"go-im/common/tcp"
)

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

type Message struct {
	Header *tcp.FixedHeader
	Body   interface{}
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnet(serverAddr),
	}

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
