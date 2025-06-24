package logic

import (
	"context"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"go_code/zhihu/application/tag/api/internal/svc"
	"go_code/zhihu/application/tag/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResourceTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetResourceTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceTagsLogic {
	return &GetResourceTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResourceTagsLogic) GetResourceTags(req *types.GetResourceTagsReq) (resp *types.GetResourceTagsResp, err error) {
	response, err := l.svcCtx.TagRpc.GetResourceTags(l.ctx, &tag.GetResourceTagsRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		return nil, err
	}

	tags := make([]types.Tag, 0, len(response.Tags))
	for _, t := range response.Tags {
		tags = append(tags, types.Tag{
			Id:         t.Id,
			TagName:    t.TagName,
			TagDesc:    t.TagDesc,
			CreateTime: t.CreateTime,
		})
	}

	return &types.GetResourceTagsResp{Tags: tags}, nil
}
