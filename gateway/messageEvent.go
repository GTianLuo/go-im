package gateway

import "go-im/common/tcp"

type MessageEvent struct {
	connId int
	Header *tcp.FixedHeader
	Body   interface{}
}

func NewMessageEvent(connId int, h *tcp.FixedHeader, body interface{}) *MessageEvent {
	return &MessageEvent{connId: connId, Header: h, Body: body}
}
