package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/internal/model"
	"time"

	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerQuestionLogic {
	return &AnswerQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerQuestionLogic) AnswerQuestion(in *qa.AnswerRequest) (*qa.AnswerReply, error) {
	a := &model.Answer{
		QuestionId: uint64(in.QuestionId),
		UserId:     uint64(in.UserId),
		Content:    in.Content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	res, err := l.svcCtx.AnswerModel.Insert(l.ctx, a)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &qa.AnswerReply{Id: id}, nil
}
