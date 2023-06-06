package middlewareConf

import (
	"fmt"
	"go-im/common/conf"
)

func GetMQUrl() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		conf.V.GetString("mq.user"),
		conf.V.GetString("mq.password"),
		conf.V.GetString("mq.ip"),
		conf.V.GetString("mq.host"))
}

func GetMQXName() string {
	return conf.V.GetString("mq.upMessage.xName")
}

func GetMQPQueueName() string {
	return conf.V.GetString("mq.upMessage.privateMsg.queueName")
}

func GetMQPBindingKey() string {
	return conf.V.GetString("mq.upMessage.privateMsg.bindingKey")
}

func GetMQGQueueName() string {
	return conf.V.GetString("mq.upMessage.groupMsg.queueName")
}

func GetMQGBindingKey() string {
	return conf.V.GetString("mq.upMessage.groupMsg.bindingKey")
}
