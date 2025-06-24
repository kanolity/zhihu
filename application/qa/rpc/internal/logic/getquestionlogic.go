package logic

import (
	"context"

	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionLogic) GetQuestion(in *qa.GetQuestionRequest) (*qa.GetQuestionReply, error) {
	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}
	return &qa.GetQuestionReply{
		Question: &qa.Question{
			Id:          int64(q.Id),
			UserId:      int64(q.UserId),
			Title:       q.Title,
			Description: q.Description,
			IsResolved:  q.IsResolved,
			CreateTime:  q.CreateTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
