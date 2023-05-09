package codec

import (
	"fmt"
	"go-im/common/tcp"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	listen, _ := net.Listen("tcp", ":9999")
	conn, _ := listen.Accept()

	codec := NewGobCodec(conn)
	h := &tcp.FixedHeader{}

	if err := codec.ReadFixedHeader(h); err != nil {
		fmt.Println(err)
		return
	}
	body1 := tcp.GetMessageBody(tcp.HeartBeatMessage)
	if err := codec.ReadBody(body1); err != nil {
		fmt.Println(err)
		return
	}
	body := tcp.GetMessageBody(tcp.PrivateChatMessage)
	if err := codec.ReadBody(body); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(h, body1, body)
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	codec := NewGobCodec(conn)
	h := &tcp.FixedHeader{Seq: 1, MessageType: tcp.PrivateChatMessage}
	body := &tcp.PrivateChat{
		From:    1,
		To:      2,
		Content: "hello world",
	}
	codec.Write(h, body)
	select {}
}
