package mq

import amqp "github.com/rabbitmq/amqp091-go"

type MQConsumer struct {
	ch *amqp.Channel
}

func NewConsumer() (*MQConsumer, error) {
	ch, err := mqConn.Channel()
	if err != nil {
		return nil, err
	}
	s := &MQConsumer{ch: ch}
	return s, nil
}

func (s *MQConsumer) ConsumeMsg(queueName string, handleMsg func(msg []byte) bool) {
	defer s.ch.Close()
	consumeCh, err := s.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	for msg := range consumeCh {
		if ok := handleMsg(msg.Body); ok {
			if err := msg.Ack(false); err != nil {
				panic(err)
			}
		}
	}
}
