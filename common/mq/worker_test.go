package mq

import (
	"go-im/common/conf"
	"go-im/common/log"
	"testing"
)

func TestMQWorker_PublishMsg(t *testing.T) {
	conf.Init("../conf/")
	if err := InitMQConn(); err != nil {
		panic(err)
	}
	worker, err := NewWorker()
	if err != nil {
		panic(err)
	}
	err = worker.PublishMsg("msg.private", []byte("sss"))
	log.Error(err)
}
