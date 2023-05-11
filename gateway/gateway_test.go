package gateway

import (
	"fmt"
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"log"
	"net"
	"testing"
)

// 测试最大连接数
func TestConn(t *testing.T) {
	for {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:4021")
		if _, err := net.DialTCP("tcp", nil, addr); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func TestSendAndRecv(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:4021")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	gobCodec := codec.NewGobCodec(conn)
	if err = gobCodec.Write(&tcp.FixedHeader{Seq: 1, MessageType: tcp.PrivateChatMessage},
		&tcp.PrivateChat{From: "122", To: "1212", Content: "hello world"}); err != nil {
		log.Fatal(err)
	}
	h := &tcp.FixedHeader{}
	if err = gobCodec.ReadFixedHeader(h); err != nil {
		log.Fatal(h)
	}
	body := tcp.GetMessageBody(h.MessageType)
	if err = gobCodec.ReadBody(body); err != nil {
		log.Fatal(err)
	}
	fmt.Println(body)
}

func TestNet(t *testing.T) {
	net.Listen("tcp", "")
}
