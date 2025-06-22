package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/article/api/internal/config"
	"go_code/zhihu/application/article/rpc/articleclient"
	"go_code/zhihu/application/user/rpc/userclient"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	ArticleRpc articleclient.Article
	UserRpc    userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	articleRPC := zrpc.MustNewClient(c.ArticleRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:     c,
		ArticleRpc: articleclient.NewArticle(articleRPC),
		UserRpc:    userclient.NewUser(userRPC),
	}
}
