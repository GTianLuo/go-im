package discovery

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"go-im/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli *clientv3.Client
	ctx context.Context
}

// NewServiceDiscovery 创建服务发现服务
func NewServiceDiscovery(ctx context.Context, endPoints []string, dailTimeOut time.Duration) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dailTimeOut,
	})
	if err != nil {
		logger.Fatal("discovery service:", err.Error())
	}
	return &ServiceDiscovery{
		cli: cli,
		ctx: ctx,
	}
}

// InitAndWatch 在服务启动时调用，初始化服务列表并开启watcher
func (dis *ServiceDiscovery) InitAndWatch(prefix string, set, del func(key, value string)) {

	//获取服务列表
	resp, err := dis.cli.Get(dis.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Fatal("discovery service:", err.Error())
	}

	//初始化服务列表
	for _, kv := range resp.Kvs {
		set(string(kv.Key), string(kv.Value))
	}
	//监听该prefix下的所有kv
	dis.watch(prefix, set, del)

}

// watch 监听指定前缀下的所有kv,直到连接关闭
func (dis *ServiceDiscovery) watch(prefix string, set func(key string, value string), del func(key string, value string)) {
	watch := dis.cli.Watch(dis.ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for resp := range watch {
		for _, event := range resp.Events {
			kv := event.Kv
			preKv := event.PrevKv
			switch event.Type {
			case mvccpb.DELETE:
				log.Info("=========================================================")
				del(string(preKv.Key), string(preKv.Value))
			case mvccpb.PUT:
				set(string(kv.Key), string(kv.Value))
			}
		}
	}
}

func (dis *ServiceDiscovery) Close() error {
	return dis.cli.Close()
}
