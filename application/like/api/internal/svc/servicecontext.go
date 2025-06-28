package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/like/api/internal/config"
	"go_code/zhihu/application/like/rpc/likeclient"
	"go_code/zhihu/application/message/rpc/messageservice"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	LikeRpc    likeclient.Like
	MessageRpc messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	likeRPC := zrpc.MustNewClient(c.LikeRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	messageRPC := zrpc.MustNewClient(c.MessageRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:     c,
		LikeRpc:    likeclient.NewLike(likeRPC),
		MessageRpc: messageservice.NewMessageService(messageRPC),
	}
}
