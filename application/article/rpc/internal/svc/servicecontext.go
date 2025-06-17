package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/article/rpc/internal/config"
	"go_code/zhihu/application/article/rpc/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(conn, c.CacheRedis),
	}
}
