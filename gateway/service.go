package gateway

import (
	"go-im/conf"
	"go-im/log"
	"net"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/conf/")
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.GetGateWayAddr())
	if err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	lis, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	if err = initManager(lis); err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	select {}
}
