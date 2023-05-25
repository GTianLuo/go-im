package serviceConf

import (
	"go-im/common/conf"
)

func GetClientLoginAddr() string {
	return conf.V.GetString("client.login")
}
