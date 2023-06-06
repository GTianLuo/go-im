package mq

import (
	"fmt"
	"go-im/common/conf"
	"go-im/common/conf/middlewareConf"
	"testing"
)

func TestMQConsumer_ConsumeMsg(t *testing.T) {
	conf.Init("../conf/")
	if err := InitMQConn(); err != nil {
		panic(err)
	}
	consumer, err := NewConsumer()
	if err != nil {
		panic(err)
	}
	consumer.ConsumeMsg(middlewareConf.GetMQPQueueName(), func(msg []byte) bool {
		fmt.Println(string(msg))
		return true
	})
}
