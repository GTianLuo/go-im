package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go-im/common/conf/middlewareConf"
)

var mqConn *amqp.Connection

func InitMQConn() error {
	conn, err := amqp.Dial(middlewareConf.GetMQUrl())
	if err != nil {
		return err
	}
	mqConn = conn
	return err
}
