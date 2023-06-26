package model

import (
	"go-im/common/message"
	"gorm.io/gorm"
)

type MsgSend struct {
	gorm.Model
	MsgFrom    string          `gorm:"index"`
	MsgTo      string          `gorm:"index"`
	MsgSeq     string          `gorm:"unique"`
	MsgContent []byte          `gorm:"column:msg_content"`
	SendTime   int64           `gorm:"index"`
	CmdType    message.CmdType `gorm:"comment: 信令消息 0: 私聊消息  1: 群聊消息"`
	MsgType    message.MsgType `gorm:"comment:消息类型 0：文本消息 1：图片消息"`
}

func (m *MsgSend) TableName() string {
	return "message_send"
}
