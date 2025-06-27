package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/follow/api/internal/config"
	"go_code/zhihu/application/follow/rpc/followclient"
	"go_code/zhihu/application/follow/rpc/types/follow"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config    config.Config
	FollowRpc follow.FollowClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	followRPC := zrpc.MustNewClient(c.FollowRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:    c,
		FollowRpc: followclient.NewFollow(followRPC),
	}
}
