package message

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"

	"testing"
)

type PrivateMsg2 struct {
	Type MsgType `protobuf:"varint,1,opt,name=Type,proto3,enum=message.pb.MsgType" json:"Type,omitempty"`
	From string  `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	To   string  `protobuf:"bytes,3,opt,name=To,proto3" json:"To,omitempty"`
	Data []byte  `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

type Cmd2 struct {
	Type    CmdType `protobuf:"varint,1,opt,name=Type,proto3,enum=message.pb.CmdType" json:"Type,omitempty"`
	Payload []byte  `protobuf:"bytes,2,opt,name=Payload,proto3" json:"Payload,omitempty"`
}

func TestProto(t *testing.T) {
	msg := &PrivateMsg{
		To:   "1",
		Type: MsgType_TextMsg,
		Data: []byte("hello world"),
	}
	bytes, _ := proto.Marshal(msg)
	c := &Cmd{
		Type:    CmdType_PrivateMsgCmd,
		Payload: bytes,
	}
	s, _ := proto.Marshal(c)

	fmt.Println(len(s))

	msg2 := &PrivateMsg2{
		From: "1",
		To:   "1",
		Type: MsgType_TextMsg,
		Data: []byte("hello world"),
	}
	s, _ = json.Marshal(msg2)
	c2 := &Cmd2{
		Type:    CmdType_PrivateMsgCmd,
		Payload: s,
	}
	s, _ = json.Marshal(c2)
	fmt.Println(len(s))

}
