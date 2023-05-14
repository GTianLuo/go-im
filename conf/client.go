package conf

func GetClientDiscoveryAddr() string {
	return v.GetString("client.discovery")
}
