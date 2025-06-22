package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/like/api/internal/config"
	"go_code/zhihu/application/like/rpc/likeclient"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config  config.Config
	LikeRpc likeclient.Like
}

func NewServiceContext(c config.Config) *ServiceContext {
	likeRPC := zrpc.MustNewClient(c.LikeRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:  c,
		LikeRpc: likeclient.NewLike(likeRPC),
	}
}
