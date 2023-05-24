package service

import (
	"context"
	"go-im/common/dao"
	"go-im/common/log"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	UnimplementedIpConfigServer
}

func InitRpcService() {
	listen, err := net.Listen("tcp", ":9998")
	if err != nil {
		panic(err)
	}
	// 创建服务
	s := grpc.NewServer()
	//注册服务
	RegisterIpConfigServer(s, &Service{})
	log.Info("rpc server on :", listen.Addr())
	//绑定服务
	if err = s.Serve(listen); err != nil {
		panic(err)
	}
}

// Auth 鉴权服务
func (s *Service) Auth(ctx context.Context, r *AuthRequest) (*AuthResponse, error) {

	response := &AuthResponse{}
	// redis 查找用户登陆状态
	userDao := dao.NewUserDao()
	userStatus, err := userDao.GetLoginStatus(r.Account)
	if err != nil {
		return response, nil
	}

	//校验token
	if r.Token != userStatus["token"] {
		return response, nil
	}
	response.Success = true
	return response, nil
}
