package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnswersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAnswersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnswersLogic {
	return &GetAnswersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnswersLogic) GetAnswers(req *types.GetAnswersReq) (resp *types.GetAnswersResp, err error) {
	response, err := l.svcCtx.QaRpc.GetAnswers(l.ctx, &qa.GetAnswersRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		Limit:      req.Limit,
	})
	if err != nil {
		return nil, err
	}

	answers := make([]types.Answer, 0, len(response.Answers))
	for _, a := range response.Answers {
		answers = append(answers, types.Answer{
			Id:         a.Id,
			QuestionId: a.QuestionId,
			UserId:     a.UserId,
			Content:    a.Content,
			CreateTime: a.CreateTime,
		})
	}

	return &types.GetAnswersResp{
		Answers: answers,
		HasMore: response.HasMore,
	}, nil
}
