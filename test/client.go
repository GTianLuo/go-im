package main

import (
	"flag"
	"fmt"
	"go-im/common/log"
	"go-im/common/tcp"
	"go-im/common/tcp/codec"
	"net"

	"time"
)

var connections = flag.Int("conn", 1, "number of tcp connections")

func main() {

	flag.Parse()
	conns := make([]codec.Codec, 0)
	for i := 0; i < *connections; i++ {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:4021")
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		conns = append(conns, codec.NewGobCodec(conn))
	}

	log.Info("完成初始化 %d 连接", len(conns))
	defer func() {
		for _, codec := range conns {
			codec.Close()
		}
	}()
	ttl := time.Second
	h := &tcp.FixedHeader{Seq: 1, MessageType: tcp.PrivateChatMessage}
	body := &tcp.PrivateChat{From: "1", To: "2", Content: "hello world"}
	for {
		for _, codec := range conns {
			if err := codec.Write(h, body); err != nil {
				log.Error(err)
				return
			}
		}
		time.Sleep(ttl)
	}
}
