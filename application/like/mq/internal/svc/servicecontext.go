package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/like/mq/internal/config"
	"go_code/zhihu/application/like/mq/internal/model"
)

type ServiceContext struct {
	Config         config.Config
	LikeModel      model.LikeModel
	LikeCountModel model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:         c,
		LikeModel:      model.NewLikeModel(conn, c.CacheRedis),
		LikeCountModel: model.NewLikeCountModel(conn, c.CacheRedis),
	}
}
