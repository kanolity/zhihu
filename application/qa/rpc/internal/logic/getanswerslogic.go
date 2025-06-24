package logic

import (
	"context"

	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnswersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAnswersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnswersLogic {
	return &GetAnswersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAnswersLogic) GetAnswers(in *qa.GetAnswersRequest) (*qa.GetAnswersReply, error) {
	list, err := l.svcCtx.AnswerModel.ListByQuestion(l.ctx, in.QuestionId, in.Cursor, in.Limit+1)
	if err != nil {
		return nil, err
	}

	resp := &qa.GetAnswersReply{
		Answers: make([]*qa.Answer, 0, in.Limit),
		HasMore: false,
	}
	for i, a := range list {
		if int64(i) < in.Limit {
			resp.Answers = append(resp.Answers, &qa.Answer{
				Id:         int64(a.Id),
				QuestionId: int64(a.QuestionId),
				UserId:     int64(a.UserId),
				Content:    a.Content,
				CreateTime: a.CreateTime.Format("2006-01-02 15:04:05"),
			})
		} else {
			resp.HasMore = true
		}
	}
	return resp, nil
}
