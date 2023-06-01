package sdk

import "go-im/common/tcp"

type Message struct {
	Header *tcp.FixedHeader
	Body   interface{}
}

func GetSystemMessage(msg string, reConn bool) *Message {
	h := &tcp.FixedHeader{
		MsgId:       -1,
		PreMsgId:    -1,
		From:        "(系统)",
		MessageType: tcp.SystemMessage,
	}
	b := &tcp.SystemMB{
		ErrMsg: msg,
		ReConn: reConn,
	}
	return &Message{h, b}
}

func GetHeartBeatMessage(from string) *Message {
	h := &tcp.FixedHeader{
		MsgId:       -1,
		PreMsgId:    -1,
		From:        from,
		MessageType: tcp.HeartBeatMessage,
	}
	b := &tcp.HeartBeatMB{}
	return &Message{h, b}
}
