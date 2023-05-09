package discovery

import (
	"context"
	"go-im/conf"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	conf.Init("../../conf/")
	discovery := NewServiceDiscovery(context.Background(), conf.GetIpConfigEndPoints(), time.Second)
	go discovery.InitAndWatch("", func(key, value string) {}, func(key, value string) {})
	time.Sleep(time.Second)
	//discovery.Close()
	time.Sleep(3 * time.Second)
}
