package codec

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

func TestProtoCodec(t *testing.T) {
	listen, _ := net.Listen("tcp", ":9999")

	conn, _ := listen.Accept()
	codec := NewProtoCodec(conn)
	for i := 0; i < 100; i++ {
		if err := codec.WriteData([]byte("你哦好" + strconv.Itoa(i))); err != nil {
			panic(err)
		}
	}
	select {}
}

func TestS(t *testing.T) {
	conn, _ := net.Dial("tcp", ":9999")
	codec := NewProtoCodec(conn)
	for {
		data, err := codec.ReadData()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
}
