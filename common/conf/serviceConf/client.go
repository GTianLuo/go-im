package serviceConf

import (
	"go-im/common/conf"
)

func GetClientLoginAddr() string {
	return conf.V.GetString("client.login")
}

func GetClientMaxReSendNums() int {
	return conf.V.GetInt("client.maxReSendCount")
}
