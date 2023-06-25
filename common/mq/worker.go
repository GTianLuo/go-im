package mq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MQWorker struct {
	ch         *amqp.Channel
	xName      string
	routingKey string
}

// NewWorker 创建一个推送消息的worker
func NewWorker(ch *amqp.Channel, xName string, routingKey string) (*MQWorker, error) {
	w := &MQWorker{ch: ch, xName: xName, routingKey: routingKey}
	return w, nil
}

// PublishMsg 推送消息
func (w *MQWorker) PublishMsg(msg []byte) error {
	return w.ch.PublishWithContext(context.Background(), w.xName, w.routingKey, false, false, amqp.Publishing{Body: msg, ContentType: "application/protobuf"})
}
