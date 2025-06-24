package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/message/api/internal/config"
	"go_code/zhihu/application/message/rpc/messageservice"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	MessageRpc messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	messageRPC := zrpc.MustNewClient(c.MessageRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:     c,
		MessageRpc: messageservice.NewMessageService(messageRPC),
	}
}
