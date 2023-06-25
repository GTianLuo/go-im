package model

/*
   历史消息表
      msgId  from  to  content
*/

type PrivateTextMsg struct {
	//gorm.Model
	Timestamp int64
	MsgId     int64
	From      string
	To        string
	Content   string
}
