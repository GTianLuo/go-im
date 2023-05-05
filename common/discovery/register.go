package discovery

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type ServiceRegister struct {
	ctx           context.Context
	cli           *clientv3.Client
	lease         clientv3.LeaseID
	key           string
	value         string
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func NewServerRegister(ctx context.Context, endPoints []string, dailTimeOut time.Duration, ttl int64, key string, value string) *ServiceRegister {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dailTimeOut,
	})
	if err != nil {
		logger.Fatal("register service", err)
	}
	s := &ServiceRegister{
		ctx:   ctx,
		cli:   cli,
		key:   key,
		value: value,
	}
	s.putKVWithLease(ttl)
	return s
}

func (s *ServiceRegister) putKVWithLease(ttl int64) {
	resp, err := s.cli.Grant(s.ctx, ttl)
	if err != nil {
		logger.Fatal("register service", err)
	}
	s.lease = resp.ID
	_, err = s.cli.Put(s.ctx, s.key, s.value, clientv3.WithLease(s.lease))
	if err != nil {
		logger.Fatal("register service", err)
	}
	s.keepAlive()
}

func (s *ServiceRegister) keepAlive() {
	keepAliveChan, err := s.cli.KeepAlive(s.ctx, s.lease)
	if err != nil {
		logger.Fatal("register service", err)
	}
	s.keepAliveChan = keepAliveChan
}

func (s *ServiceRegister) ListenKeepAliveChan() {
	for _ = range s.keepAliveChan {
		logger.Info(s.key, "successful keep alive")
	}
	logger.Info(s.key, "failed to keep alive")
}

func (s *ServiceRegister) Close() {
	defer s.cli.Close()
	//撤销租约
	_, _ = s.cli.Revoke(s.ctx, s.lease)
	logger.Info(s.key, "has revoke lease")
}
