package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/tag/rpc/internal/config"
	"go_code/zhihu/application/tag/rpc/internal/model"
)

type ServiceContext struct {
	Config           config.Config
	TagModel         model.TagModel
	TagResourceModel model.TagResourceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:           c,
		TagModel:         model.NewTagModel(conn),
		TagResourceModel: model.NewTagResourceModel(conn),
	}
}
