package mq

import (
	"go-im/common/conf"
	"go-im/common/conf/middlewareConf"
	"go-im/common/conf/serviceConf"
	"testing"
)

func TestMQWorker_PublishMsg(t *testing.T) {
	conf.Init("../conf/")
	if err := middlewareConf.InitMQConn(); err != nil {
		panic(err)
	}
	if err := middlewareConf.InitMqService(); err != nil {
		panic(err)
	}
	ch, err := middlewareConf.GetMqChannel()
	if err != nil {
		panic(err)
	}
	err = QueueBind(ch, serviceConf.GetGateWayMqXName(), serviceConf.GetGatewayMqQueueName())
	if err != nil {
		panic(err)
	}

	worker, err := NewWorker(ch, serviceConf.GetGateWayMqXName(), serviceConf.GetGatewayMqQueueName())
	if err != nil {
		panic(err)
	}
	if err = worker.PublishMsg([]byte("hello world")); err != nil {
		panic(err)
	}
}
