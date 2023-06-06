package gateway

import (
	"go-im/common/conf"
	"go-im/common/conf/dbConf"
	"go-im/common/log"
	"go-im/common/timingwheel"
	"go-im/gateway/rpc/client"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/common/conf/")
	// 初始化db
	dbConf.InitDbService()
	client.RpcClientInit()
	timingwheel.InitTimingWheel()
	if err := initManager(); err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	select {}
}
