package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/internal/model"
	"time"

	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateQuestionLogic) CreateQuestion(in *qa.CreateQuestionRequest) (*qa.CreateQuestionReply, error) {
	now := time.Now()
	q := &model.Question{
		UserId:      uint64(in.UserId),
		Title:       in.Title,
		Description: in.Description,
		IsResolved:  false,
		CreateTime:  now,
		UpdateTime:  now,
	}
	res, err := l.svcCtx.QuestionModel.Insert(l.ctx, q)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &qa.CreateQuestionReply{Id: id}, nil
}
