package sdk

const (
	MsgTypeText = "text"
)

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

type Message struct {
	Type       string
	Name       string
	FormUserID string
	ToUserID   string
	Content    string
	Session    string
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnet(serverAddr),
	}
}
func (chat *Chat) Send(msg *Message) {
	chat.conn.send(msg)
}

//Close close chat
func (chat *Chat) Close() {
	chat.conn.close()
}

//Recv receive message
func (chat *Chat) Recv() <-chan *Message {
	return chat.conn.recv()
}
