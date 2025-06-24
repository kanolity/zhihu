package logic

import (
	"context"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"go_code/zhihu/application/tag/api/internal/svc"
	"go_code/zhihu/application/tag/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagReq) (resp *types.CreateTagResp, err error) {
	response, err := l.svcCtx.TagRpc.CreateTag(l.ctx, &tag.CreateTagRequest{
		TagName: req.TagName,
		TagDesc: req.TagDesc,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateTagResp{Id: response.Id}, nil
}
