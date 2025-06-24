package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionLogic) CreateQuestion(req *types.CreateQuestionReq) (resp *types.CreateQuestionResp, err error) {
	response, err := l.svcCtx.QaRpc.CreateQuestion(l.ctx, &qa.CreateQuestionRequest{
		UserId:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateQuestionResp{Id: response.Id}, nil
}
