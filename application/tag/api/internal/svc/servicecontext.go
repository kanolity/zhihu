package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/tag/api/internal/config"
	"go_code/zhihu/application/tag/rpc/tagservice"
	"go_code/zhihu/pkg/interceptors"
)

type ServiceContext struct {
	Config config.Config
	TagRpc tagservice.TagService
}

func NewServiceContext(c config.Config) *ServiceContext {
	tagRPC := zrpc.MustNewClient(c.TagRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config: c,
		TagRpc: tagservice.NewTagService(tagRPC),
	}
}
