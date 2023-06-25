package middlewareConf

import (
	"go-im/common/conf"
	"go-im/common/conf/serviceConf"
	"go-im/common/mq"
	"testing"
)

func TestMq(t *testing.T) {

	conf.Init("../")
	if err := InitMQConn(); err != nil {
		panic(err)
	}
	if err := InitMqService(); err != nil {
		panic(err)
	}
	ch, err := GetMqChannel()
	if err != nil {
		panic(err)
	}
	err = QueueBind(serviceConf.GetGateWayMqXName(), serviceConf.GetGatewayMqQueueName())
	if err != nil {
		panic(err)
	}

	worker, err := mq.NewWorker(ch, serviceConf.GetGateWayMqXName(), serviceConf.GetGatewayMqQueueName())
	if err != nil {
		panic(err)
	}
	if err = worker.PublishMsg([]byte("hello world")); err != nil {
		panic(err)
	}
}
