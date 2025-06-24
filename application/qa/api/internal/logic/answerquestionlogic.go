package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAnswerQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerQuestionLogic {
	return &AnswerQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AnswerQuestionLogic) AnswerQuestion(req *types.AnswerReq) (resp *types.AnswerResp, err error) {
	response, err := l.svcCtx.QaRpc.AnswerQuestion(l.ctx, &qa.AnswerRequest{
		QuestionId: req.QuestionId,
		UserId:     req.UserId,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &types.AnswerResp{Id: response.Id}, nil
}
