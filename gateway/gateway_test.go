package gateway

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"go-im/common/log"
	"net"
	"testing"
)

// 测试最大连接数
func TestConn(t *testing.T) {
	conns := make([]*net.TCPConn, 0)
	for i := 0; i < 20000; i++ {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:4021")
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		conns = append(conns, conn)

	}
	select {}
}

func TestWorkPool(t *testing.T) {
	pool, _ := ants.NewPool(10)
	for i := 0; i < 100; i++ {
		pool.Submit(func() {

		})
	}
}

func TestChannel(t *testing.T) {
	ch := make(chan bool, 1)
	close(ch)
	select {
	case c := <-ch:
		log.Info(c)
	default:
		log.Info(222)
	}
}

func TestIO(t *testing.T) {
	pool, _ := ants.NewPool(10)
	defer pool.Free()

	for i := 0; i < 100; i++ {
		taskID := i
		pool.Submit(func() {
			fmt.Println(taskID)
		})
	}
}
