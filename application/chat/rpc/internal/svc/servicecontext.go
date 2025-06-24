package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/chat/rpc/internal/config"
	"go_code/zhihu/application/chat/rpc/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	SessionModel model.ChatSessionModel
	MessageModel model.ChatMessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:       c,
		SessionModel: model.NewChatSessionModel(conn),
	}
}
