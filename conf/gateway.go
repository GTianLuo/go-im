package conf

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
