package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/message/rpc/internal/config"
	"go_code/zhihu/application/message/rpc/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel model.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:       c,
		MessageModel: model.NewMessageModel(conn),
	}
}
