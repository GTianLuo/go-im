package client

import (
	"context"
	"go-im/common/conf/serviceConf"
	"go-im/common/log"
	"go-im/ipConfig/rpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var ipConfigClient service.IpConfigClient

func initIpConfigClient() {
	conn, err := grpc.Dial(serviceConf.GetGateWayAuthAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = conn.Close()
		panic(err)
	}
	ipConfigClient = service.NewIpConfigClient(conn)
}

func Auth(account, token string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10000)
	resp, err := ipConfigClient.Auth(ctx, &service.AuthRequest{Account: account, Token: token})
	if err != nil {
		log.Error(err)
		return false, err
	}
	return resp.Success, err
}
