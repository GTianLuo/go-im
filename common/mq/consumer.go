package mq

import (
	"github.com/panjf2000/ants/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"go-im/common/log"
)

type MQConsumer struct {
	ch        *amqp.Channel
	queueName string //要消费的队列名
	workPool  *ants.Pool
}

func NewConsumer(ch *amqp.Channel, queueName string, workPool *ants.Pool) (*MQConsumer, error) {
	s := &MQConsumer{ch: ch, queueName: queueName, workPool: workPool}
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
			for true {
				msgBody := msg.Body

				if err := s.workPool.Submit(
					func() {
						var errAck error
						if handleMsg(msgBody) {
							errAck = msg.Ack(false)
						}
						if errAck != nil {
							log.Error(errAck)
						}
					},
				); err == nil {
					break
				}
			}
		}
	}

}
