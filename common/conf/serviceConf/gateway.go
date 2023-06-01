package serviceConf

import (
	"go-im/common/conf"
	"time"
)

func GetGateWayMaxConnsNum() int32 {
	return conf.V.GetInt32("gateway.maxConnsNum")
}

func GetGateWayWorkPoolNum() int {
	return conf.V.GetInt("gateway.workPoolNum")
}

func GetGateWayDeviceId() int32 {
	return conf.V.GetInt32("gateway.deviceId")
}

func GetGateWayAddr() string {
	return conf.V.GetString("gateway.addr")
}

func GetGatewayEndPoints() []string {
	strings := conf.V.Get("gateway.endPoints")
	endPoints := make([]string, 1)
	for _, s := range strings.([]interface{}) {
		endPoints = append(endPoints, s.(string))
	}
	return endPoints
}

func GetGatewayLeaseDDL() int64 {
	return conf.V.GetInt64("gateway.leaseDDL")
}

func GetGatewayDailTimeOut() time.Duration {
	return time.Duration(conf.V.Get("gateway.dailTimeOut").(int)) * time.Second
}

func GetGateWayEpollMaxTriggerConn() int {
	return conf.V.GetInt("gateway.epoll.maxTriggerEvent")
}

func GetGateWayReactorNums() int {
	return conf.V.GetInt("gateway.reactorNum")
}

func GetGateWayAuthAddr() string {
	return conf.V.GetString("gateway.authAddr")
}

func GetGateWayHeartbeatTimeout() time.Duration {
	return conf.V.GetDuration("gateway.heartBeatTimeout") * time.Second
}
