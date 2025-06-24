package logic

import (
	"context"

	"go_code/zhihu/application/tag/rpc/internal/svc"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResourceTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetResourceTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceTagsLogic {
	return &GetResourceTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetResourceTagsLogic) GetResourceTags(in *tag.GetResourceTagsRequest) (*tag.GetResourceTagsReply, error) {
	relations, err := l.svcCtx.TagResourceModel.FindByBizTarget(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return &tag.GetResourceTagsReply{Tags: []*tag.Tag{}}, nil
	}

	tagIds := make([]int64, 0, len(relations))
	for _, r := range relations {
		tagIds = append(tagIds, int64(r.TagId))
	}

	tags, err := l.svcCtx.TagModel.BatchGetTags(l.ctx, tagIds)
	if err != nil {
		return nil, err
	}

	respTags := make([]*tag.Tag, 0, len(tags))
	for _, t := range tags {
		respTags = append(respTags, &tag.Tag{
			Id:         int64(t.Id),
			TagName:    t.TagName,
			TagDesc:    t.TagDesc,
			CreateTime: t.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return &tag.GetResourceTagsReply{Tags: respTags}, nil
}
