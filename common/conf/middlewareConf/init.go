package middlewareConf

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

var mqConn *amqp.Connection

func InitMQConn() error {
	conn, err := amqp.Dial(GetMQUrl())
	if err != nil {
		return err
	}
	mqConn = conn
	return err
}

// InitMqService 初始化mq服务，创建exchange
func InitMqService() error {
	if mqConn == nil {
		return errors.New("must init mq connection first")
	}
	ch, err := mqConn.Channel()
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(GetMQXName(), "direct", true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	return nil
}

// QueueBind 该方法负责声明queue，并将queue和exchange绑定
// 绑定队列和mq时，使用queueName作为routingKey
func QueueBind(xName string, queueName string) error {
	ch, err := mqConn.Channel()
	if err != nil {
		panic(err)
	}
	pQueue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	// 绑定privateQueue和交换机
	if err = ch.QueueBind(pQueue.Name, queueName, xName, false, nil); err != nil {
		return err
	}
	return nil
}

func GetMqChannel() (*amqp.Channel, error) {
	return mqConn.Channel()
}
