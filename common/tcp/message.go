package tcp

type MessageType int8

const (
	ErrorMessage MessageType = iota
	HeartBeatMessage
	PrivateChatMessage
	GroupChatMessage
	AuthMessage
)

// FixedHeader 固定头部
type FixedHeader struct {
	To          int
	Seq         int64       //序号 4个字节
	MessageType MessageType //消息类型 1个字节
}

type PrivateChat struct {
	From    int
	To      int
	Content string
}

type GroupChat struct {
	From  int
	Group int
	//Content string
}

type HeartBeat struct {
}

func GetMessageBody(t MessageType) interface{} {
	switch t {
	case PrivateChatMessage:
		return new(PrivateChat)
	case HeartBeatMessage:
		return new(HeartBeat)
	case GroupChatMessage:
		return new(GroupChat)
	}
	return nil
}
