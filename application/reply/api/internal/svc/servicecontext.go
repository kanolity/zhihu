package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/reply/api/internal/config"
	"go_code/zhihu/application/reply/rpc/replyservice"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config   config.Config
	ReplyRpc replyservice.ReplyService
}

func NewServiceContext(c config.Config) *ServiceContext {
	replyRPC := zrpc.MustNewClient(c.ReplyRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:   c,
		ReplyRpc: replyservice.NewReplyService(replyRPC),
	}
}
