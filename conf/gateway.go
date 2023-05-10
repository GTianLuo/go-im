package conf

import "time"

func GetGateWayMaxConnsNum() int32 {
	return v.GetInt32("gateway.maxConnsNum")
}

func GetGateWayWorkPoolNum() int {
	return v.GetInt("gateway.workPoolNum")
}

func GetGateWayDeviceId() int32 {
	return v.GetInt32("gateway.deviceId")
}

func GetGateWayAddr() string {
	return v.GetString("gateway.addr")
}

func GetGatewayEndPoints() []string {
	strings := v.Get("gateway.endPoints")
	endPoints := make([]string, 1)
	for _, s := range strings.([]interface{}) {
		endPoints = append(endPoints, s.(string))
	}
	return endPoints
}

func GetGatewayLeaseDDL() int64 {
	return v.GetInt64("gateway.leaseDDL")
}

func GetGatewayDailTimeOut() time.Duration {
	return time.Duration(v.Get("gateway.dailTimeOut").(int)) * time.Second
}
