package util

import (
	"context"
	"go-im/common/conf/serviceConf"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"sync"
)

func NewLock() sync.Locker {
	cli, err := clientv3.New(clientv3.Config{Endpoints: serviceConf.GetIpConfigEndPoints()})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	session, err := concurrency.NewSession(cli, concurrency.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	return concurrency.NewLocker(session, "/myLock/")
}
