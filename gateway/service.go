package gateway

import (
	"go-im/common/conf"
	"go-im/common/conf/dbConf"
	"go-im/common/conf/middlewareConf"
	"go-im/common/conf/serviceConf"
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
	// 初始化mq连接
	if err := middlewareConf.InitMQConn(); err != nil {
		panic(err)
	}
	if err := middlewareConf.InitMqService(); err != nil {
		panic(err)
	}
	// 初始化当前gateway与mq关联的queue
	if err := middlewareConf.QueueBind(serviceConf.GetGateWayMqXName(), serviceConf.GetGatewayMqQueueName()); err != nil {
		panic(err)
	}
	timingwheel.InitTimingWheel()
	if err := initManager(); err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	select {}
}
