package source

import (
	"github.com/bytedance/gopkg/util/logger"
	"go-im/common/discovery"
)

type EventType int

const (
	AddNodeEvent EventType = iota
	DelNodeEvent
)

var eventChan chan *Event

type Event struct {
	Type         EventType
	Ip           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

func NewEvent(key, value string) *Event {
	//解析value
	ed := &discovery.EndpointInfo{}
	if err := ed.UnMarshal([]byte(value)); err != nil {
		logger.Fatal("failed to unmrshal:", err)
	}
	connectNum, ok := ed.Metadata["connect_num"].(float64)
	if !ok {
		panic("key does not exist in the metadata")
	}
	messageBytes, ok := ed.Metadata["message_bytes"].(float64)
	if !ok {
		panic("key does not exist in the metadata")
	}
	return &Event{
		Port:         ed.Port,
		Ip:           ed.IP,
		ConnectNum:   connectNum,
		MessageBytes: messageBytes,
	}
}

func EventChan() <-chan *Event {
	return eventChan
}

func (e *Event) Key() string {
	return e.Ip + ":" + e.Port
}
