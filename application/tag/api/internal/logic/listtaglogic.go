package logic

import (
	"context"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"go_code/zhihu/application/tag/api/internal/svc"
	"go_code/zhihu/application/tag/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTagLogic {
	return &ListTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTagLogic) ListTag(req *types.ListTagReq) (resp *types.ListTagResp, err error) {
	response, err := l.svcCtx.TagRpc.ListTag(l.ctx, &tag.ListTagRequest{
		Cursor: req.Cursor,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.Tag, 0, len(response.Tags))
	for _, t := range response.Tags {
		list = append(list, types.Tag{
			Id:         t.Id,
			TagName:    t.TagName,
			TagDesc:    t.TagDesc,
			CreateTime: t.CreateTime,
		})
	}

	return &types.ListTagResp{
		Tags:    list,
		HasMore: response.HasMore,
	}, nil
}
