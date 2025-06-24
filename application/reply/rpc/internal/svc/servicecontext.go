package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/reply/rpc/internal/config"
	"go_code/zhihu/application/reply/rpc/internal/model"
)

type ServiceContext struct {
	Config          config.Config
	ReplyModel      model.ReplyModel
	ReplyCountModel model.ReplyCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		ReplyModel:      model.NewReplyModel(conn),
		ReplyCountModel: model.NewReplyCountModel(conn),
	}
}
