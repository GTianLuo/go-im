package gateway

import (
	"go-im/common/proto/message"
)

type MessageEvent struct {
	connId int
	cmd    *message.Cmd
}

func NewMessageEvent(connId int, cmd *message.Cmd) *MessageEvent {
	return &MessageEvent{connId: connId, cmd: cmd}
}
