package tcp

type MessageType int8

const (
	ErrorMessage MessageType = iota
	HeartBeatMessage
	PrivateChatMessage
	GroupChatMessage
	AuthMessage
	AuthResponseMessage
	SystemMessage
	AckMessage
)

type AckType int8

// FixedHeader 固定头部
type FixedHeader struct {
	MsgId       int64
	From        string
	MessageType MessageType //消息类型 1个字节
}

// PrivateChatMB  私人聊天消息
type PrivateChatMB struct {
	To      string
	Content string
}

// GroupChatMB 群聊消息
type GroupChatMB struct {
	Group   string
	Content string
}

// AuthMB 登陆鉴权消息
type AuthMB struct {
	Token string
}

// AuthResponseMB 登陆鉴权响应
type AuthResponseMB struct {
	Status int    // -1表示内部错误，0表示鉴权失败，1表示成功,2表示该网关连接达到上限，需重新连接
	ErrMsg string //错误信息
}

// HeartBeatMB 心跳消息，空结构体
type HeartBeatMB struct {
}

// SystemMB 系统消息
type SystemMB struct {
	ErrMsg string
	ReConn bool // 是否需要重连
}

// MessageAck ack
type AckMB struct {
	//AckType AckType
}

func GetMessageBody(t MessageType) interface{} {
	switch t {
	case PrivateChatMessage:
		return new(PrivateChatMB)
	case HeartBeatMessage:
		return new(HeartBeatMB)
	case GroupChatMessage:
		return new(GroupChatMB)
	case AuthMessage:
		return new(AuthMB)
	case AckMessage:
		return new(AckMB)
	case AuthResponseMessage:
	}
	return nil
}
