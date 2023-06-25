package mq

import amqp "github.com/rabbitmq/amqp091-go"

type MQConsumer struct {
	ch        *amqp.Channel
	queueName string //要消费的队列名
}

func NewConsumer(ch *amqp.Channel, queueName string) (*MQConsumer, error) {
	s := &MQConsumer{ch: ch, queueName: queueName}
	return s, nil
}

// ConsumeMsg 从消息队列中消费消息，该方法会一直占有协程资源(循环)
// done :控制循环任务的关闭
// handleMsg : 回调函数，获取消息后执行该回调函数，执行成功回调函数返回true，并发送ack
func (s *MQConsumer) ConsumeMsg(done chan struct{}, handleMsg func(msg []byte) bool) {
	defer s.ch.Close()
	consumeCh, err := s.ch.Consume(s.queueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-done:
			return
		case msg := <-consumeCh:
			if ok := handleMsg(msg.Body); ok {
				if err := msg.Ack(false); err != nil {
					panic(err)
				}
			}
		}
	}

}
