package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestAuth(t *testing.T) {
	conn, err := grpc.Dial(":9998", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := NewIpConfigClient(conn)
	auth, err := client.Auth(context.Background(), &AuthRequest{Account: "2985496686", Token: "e8d7aab4-22fc-4d3a-8a0e-5e63ef0ce77"})
	fmt.Println(auth.Success)
}

func TestS(t *testing.T) {
	fmt.Println(`{"account":` + `"` + "account" + `",` + `"password":` + `"` + "password" + `"}`)
}
