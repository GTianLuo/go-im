package serviceConf

import (
	"go-im/common/conf"
	"time"
)

func GetIpConfigEndPoints() []string {
	strings := conf.V.Get("ipConfig.endPoints")
	endPoints := make([]string, 1)
	for _, s := range strings.([]interface{}) {
		endPoints = append(endPoints, s.(string))
	}
	return endPoints
}

func GetIpConfigEtcdServer() string {
	return conf.V.Get("ipConfig.etcdService").(string)
}

func GetIpConfigDailTimeOut() time.Duration {
	return time.Duration(conf.V.Get("ipConfig.dailTimeOut").(int)) * time.Second
}
