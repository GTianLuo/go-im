package source

import (
	"context"
	"go-im/common/discovery"
	"go-im/conf"
)

func Init() {
	eventChan = make(chan *Event, 5)
	ctx := context.Background()
	go handleData(ctx)
}

func handleData(ctx context.Context) {

	//创建服务发现对象
	d := discovery.NewServiceDiscovery(ctx, conf.GetIpConfigEndPoints(), conf.GetIpConfigDailTimeOut())

	defer d.Close()
	//创建闭包函数，用户服务变更后处理
	addFunc := func(key, value string) {
		event := NewEvent(key, value)
		event.Type = AddNodeEvent
		eventChan <- event
	}

	delFunc := func(key, value string) {
		event := NewEvent(key, value)
		event.Type = DelNodeEvent
		eventChan <- event
	}

	//初始化IpConfig服务数据并监听
	d.InitAndWatch(conf.GetIpConfigEtcdServer(), addFunc, delFunc)
}
