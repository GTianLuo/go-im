package gateway

import "go-im/common/tcp"

type MessageEvent struct {
	UserID int
	Header *tcp.FixedHeader
	Body   interface{}
}

func NewMessageEvent(userId int, h *tcp.FixedHeader, body interface{}) *MessageEvent {
	return &MessageEvent{UserID: userId, Header: h, Body: body}
}
