package logic

import (
	"context"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"go_code/zhihu/application/tag/api/internal/svc"
	"go_code/zhihu/application/tag/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTagToResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTagToResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagToResourceLogic {
	return &AddTagToResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTagToResourceLogic) AddTagToResource(req *types.AddTagToResourceReq) (resp *types.AddTagToResourceResp, err error) {
	_, err = l.svcCtx.TagRpc.AddTagToResource(l.ctx, &tag.AddTagToResourceRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		TagId:    req.TagId,
		UserId:   req.UserId,
	})
	return &types.AddTagToResourceResp{}, err
}
