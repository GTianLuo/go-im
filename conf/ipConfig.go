package conf

import (
	"time"
)

func GetIpConfigEndPoints() []string {
	strings := v.Get("ipConfig.endPoints")
	endPoints := make([]string, 1)
	for _, s := range strings.([]interface{}) {
		endPoints = append(endPoints, s.(string))
	}
	return endPoints
}

func GetIpConfigEtcdServer() string {
	return v.Get("ipConfig.etcdService").(string)
}

func GetIpConfigDailTimeOut() time.Duration {
	return time.Duration(v.Get("ipConfig.dailTimeOut").(int)) * time.Second
}

func GetIpConfigLeaseDDL() int64 {
	return v.GetInt64("ipConfig.leaseDDL")
}
