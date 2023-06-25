package messaging

import (
	"go-im/common/conf"
	"go-im/common/conf/dbConf"
	"go-im/common/conf/middlewareConf"
	"go-im/common/conf/serviceConf"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/common/conf/")
	//初始化mq
	if err := middlewareConf.InitMQConn(); err != nil {
		panic(err)
	}
	// 初始化mq对应的queue
	if err := middlewareConf.QueueBind(serviceConf.GetMessagingMqXName(), serviceConf.GetMessagingMqQueueName()); err != nil {
		panic(err)
	}
	dbConf.InitDbService()
	if err := initMsgHandle(); err != nil {
		panic(err)
	}
	select {}
}
