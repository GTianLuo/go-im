package sdk

import (
	"go-im/common/proto/message"
)

func GetSystemMessage(msg string, reConn bool) *message.Cmd {
	cmd := &message.Cmd{
		Type:    message.CmdType_SystemCmd,
		Payload: []byte(msg),
		From:    "用户",
	}
	if reConn {
		cmd.MsgId = 1
	} else {
		cmd.MsgId = 0
	}
	return nil
}

func GetHeartBeatMessage(from string) *message.Cmd {
	cmd := &message.Cmd{
		Type: message.CmdType_HeartBeatCmd,
		From: from,
	}
	return cmd
}
