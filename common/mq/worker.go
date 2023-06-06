package mq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go-im/common/conf/middlewareConf"
)

type MQWorker struct {
	ch *amqp.Channel
}

// 用于初始化MQ
func (w *MQWorker) initUpMsgExchangeQueue() error {
	ch := w.ch
	err := ch.ExchangeDeclare(middlewareConf.GetMQXName(), "direct", true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	// 声明处理privateMsg 的queue
	pQueue, err := ch.QueueDeclare(middlewareConf.GetMQPQueueName(), true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	// 绑定privateQueue和交换机
	err = ch.QueueBind(pQueue.Name, middlewareConf.GetMQPBindingKey(), middlewareConf.GetMQXName(), false, nil)
	//声明处理groupMsg 的queue
	gQueue, err := ch.QueueDeclare(middlewareConf.GetMQGQueueName(), true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	// 绑定groupQueue和交换机
	err = ch.QueueBind(gQueue.Name, middlewareConf.GetMQGBindingKey(), middlewareConf.GetMQXName(), false, nil)
	return nil
}

func NewWorker() (*MQWorker, error) {
	ch, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}
	w := &MQWorker{ch: ch}
	/*
		if err := w.initUpMsgExchangeQueue(); err != nil{
			panic(err)
		}*/
	return w, nil
}

func (w *MQWorker) PublishMsg(routingKey string, msg []byte) error {
	return w.ch.PublishWithContext(context.Background(), middlewareConf.GetMQXName(), routingKey, false, false, amqp.Publishing{Body: msg, ContentType: "application/protobuf"})
}
