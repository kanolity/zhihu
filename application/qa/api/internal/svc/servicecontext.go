package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/qa/api/internal/config"
	"go_code/zhihu/application/qa/rpc/qaservice"
	"go_code/zhihu/application/user/rpc/userclient"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config  config.Config
	QaRpc   qaservice.QaService
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	qaRPC := zrpc.MustNewClient(c.QaRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:  c,
		QaRpc:   qaservice.NewQaService(qaRPC),
		UserRpc: userclient.NewUser(userRPC),
	}
}
