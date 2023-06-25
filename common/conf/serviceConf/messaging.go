package serviceConf

import (
	"go-im/common/conf"
	"strconv"
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

func GetMessagingMqRoutingKey(deviceId int) string {
	return conf.V.GetString("messaging.mq.routingKeyPrefix") + strconv.Itoa(deviceId)
}
