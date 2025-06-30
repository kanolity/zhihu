package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/article/rpc/articleclient"
	"go_code/zhihu/application/message/rpc/messageservice"
	"go_code/zhihu/application/reply/api/internal/config"
	"go_code/zhihu/application/reply/rpc/replyservice"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	ReplyRpc   replyservice.ReplyService
	MessageRpc messageservice.MessageService
	ArticleRpc articleclient.Article
}

func NewServiceContext(c config.Config) *ServiceContext {
	replyRPC := zrpc.MustNewClient(c.ReplyRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	messageRPC := zrpc.MustNewClient(c.MessageRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	articleRPC := zrpc.MustNewClient(c.ArticleRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:     c,
		ReplyRpc:   replyservice.NewReplyService(replyRPC),
		MessageRpc: messageservice.NewMessageService(messageRPC),
		ArticleRpc: articleclient.NewArticle(articleRPC),
	}
}
