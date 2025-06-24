package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/qa/rpc/internal/config"
	"go_code/zhihu/application/qa/rpc/internal/model"
)

type ServiceContext struct {
	Config        config.Config
	QuestionModel model.QuestionModel
	AnswerModel   model.AnswerModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:        c,
		QuestionModel: model.NewQuestionModel(conn),
		AnswerModel:   model.NewAnswerModel(conn),
	}
}
