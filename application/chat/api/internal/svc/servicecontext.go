package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/chat/api/internal/config"
	"go_code/zhihu/application/chat/rpc/chatservice"
	"go_code/zhihu/application/user/rpc/userclient"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config  config.Config
	ChatRpc chatservice.ChatService
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	chatRPC := zrpc.MustNewClient(c.ChatRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:  c,
		ChatRpc: chatservice.NewChatService(chatRPC),
		UserRpc: userclient.NewUser(userRPC),
	}
}
