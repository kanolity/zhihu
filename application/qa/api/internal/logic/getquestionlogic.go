package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionLogic) GetQuestion(req *types.GetQuestionReq) (resp *types.GetQuestionResp, err error) {
	response, err := l.svcCtx.QaRpc.GetQuestion(l.ctx, &qa.GetQuestionRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	q := response.Question
	return &types.GetQuestionResp{
		Question: types.Question{
			Id:          q.Id,
			UserId:      q.UserId,
			Title:       q.Title,
			Description: q.Description,
			IsResolved:  q.IsResolved,
			CreateTime:  q.CreateTime,
		},
	}, nil
}
