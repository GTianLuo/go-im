package serviceConf

import (
	"go-im/common/conf"
)

func GetClientDiscoveryAddr() string {
	return conf.V.GetString("client.discovery")
}
