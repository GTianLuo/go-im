package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cast"
	"go-im/common/conf"
	serviceConf2 "go-im/common/conf/serviceConf"
	"go-im/common/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	// 获取 CPU 占用率
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(percent)
}

func TestSort(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	fmt.Println(s)
}

func TestLock(t *testing.T) {
	conf.Init("../../conf/")
	locker := NewLock()
	locker.Lock()
	locker.Unlock()
}

func TestVersion(t *testing.T) {
	conf.Init("../../conf/")
	cli, err := clientv3.New(clientv3.Config{Endpoints: serviceConf2.GetIpConfigEndPoints()})
	if err != nil {
		log.Fatal(err)
	}
	//getOwner是通过前缀来范围查找，WithFirstCreate()筛选出当前存在的最小revision对应的值
	getOwner := clientv3.OpGet("lock/", clientv3.WithLastCreate()...)
	response, err := cli.Do(context.Background(), getOwner)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Get().Header.Revision)
	fmt.Println(response.Get().Kvs)
}

func TestWork(t *testing.T) {
	conf.Init("../../conf/")
	cli, err := clientv3.New(clientv3.Config{Endpoints: serviceConf2.GetGatewayEndPoints()})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := concurrency.NewSTM(cli, func(stm concurrency.STM) error {
		p1 := cast.ToInt(stm.Get("p1"))
		_ = cast.ToInt(stm.Get("p2"))
		if p1 < 0 {
			return errors.New("too little money")
		}
		stm.Put("p1", cast.ToString(p1-100))
		stm.Put("p2", "300")
		return nil
	})
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(resp.Succeeded)
}

func TransferMoney(cli *clientv3.Client, from string, to string, account int) {
	resp, err := concurrency.NewSTM(cli, func(stm concurrency.STM) error {
		from := cast.ToInt(stm.Get(from))
		to := cast.ToInt(stm.Get(to))
		if from < 0 {
			//事务回滚
			return errors.New("too little money")
		}
		stm.Put("p1", cast.ToString(to-account))
		stm.Put("p2", cast.ToString(from+account))
		return nil
	})
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(resp.Succeeded)
}

func TestTean(t *testing.T) {
	conf.Init("../../conf/")
	cli, err := clientv3.New(clientv3.Config{Endpoints: serviceConf2.GetGatewayEndPoints()})
	if err != nil {
		log.Fatal(err)
	}
	txn := cli.Txn(context.Background())
	txn.Commit()
}

func TestLock2(t *testing.T) {
	var mu sync.RWMutex
	go func() {
		mu.RLock()
		fmt.Println("111111111")
		time.Sleep(time.Second * 20)
		mu.RUnlock()
	}()
	time.Sleep(time.Second)
	go func() {
		mu.Lock()
		fmt.Println("2222222222222222")
		mu.Unlock()
	}()
	time.Sleep(time.Second)
	go func() {
		mu.RLock()
		fmt.Println("333333333333")
		mu.RUnlock()
	}()
	time.Sleep(time.Second * 1000)
}

func TestG(t *testing.T) {
	for i := 0; i < 4000000; i++ {
		go func() {
			for true {
				time.Sleep(time.Second)
				fmt.Println(".")
			}
		}()
	}
	time.Sleep(1000000 * time.Second)
}
