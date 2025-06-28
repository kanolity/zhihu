package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/article/api/internal/config"
	"go_code/zhihu/application/article/rpc/articleclient"
	"go_code/zhihu/application/message/rpc/messageservice"
	"go_code/zhihu/application/tag/rpc/tagservice"
	"go_code/zhihu/application/user/rpc/userclient"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	ArticleRpc articleclient.Article
	UserRpc    userclient.User
	TagRpc     tagservice.TagService
	MessageRpc messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	articleRPC := zrpc.MustNewClient(c.ArticleRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	tagRPC := zrpc.MustNewClient(c.TagRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	messageRPC := zrpc.MustNewClient(c.MessageRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:     c,
		ArticleRpc: articleclient.NewArticle(articleRPC),
		UserRpc:    userclient.NewUser(userRPC),
		TagRpc:     tagservice.NewTagService(tagRPC),
		MessageRpc: messageservice.NewMessageService(messageRPC),
	}
}
