package serviceConf

import (
	"go-im/common/conf"
)

func GetMessagingMqConsumerNum() int {
	return conf.V.GetInt("messaging.mqConsumerNum")
}

func GetMessagingMqQueueName() string {
	return conf.V.GetString("messaging.mq.queueName")
}

func GetMessagingMqXName() string {
	return conf.V.GetString("messaging.mq.xName")
}

func GetMessagingMqRoutingKey(deviceId string) string {
	return conf.V.GetString("messaging.mq.routingKeyPrefix") + deviceId
}
