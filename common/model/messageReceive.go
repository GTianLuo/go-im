package model

import "gorm.io/gorm"

type MsgReceive struct {
	gorm.Model
	MsgFrom string `gorm:"index"`
	MsgTo   string `gorm:"index"`
	MsgSeq  string
}

func (m *MsgReceive) TableName() string {
	return "message_receive"
}
