package client

import (
	"context"
	"fmt"
	"go-im/common/conf/serviceConf"
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
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	resp, err := ipConfigClient.Auth(ctx, &service.AuthRequest{Account: account, Token: token})
	fmt.Println("=========================")
	if err != nil {
		return false, err
	}
	return resp.Success, err
}