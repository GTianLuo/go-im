package conf

import "github.com/spf13/viper"

func GetGateWayMaxConnsNum() int32 {
	return viper.GetInt32("gateway.maxConnsNum")
}

func GetGateWayWorkPoolNum() int {
	return viper.GetInt("gateway.workPoolNum")
}

func GetGateWayDeviceId() int32 {
	return viper.GetInt32("gateway.deviceId")
}

func GetGateWayAddr() string {
	return viper.GetString("gateway.addr")
}
