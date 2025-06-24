package logic

import (
	"context"

	"go_code/zhihu/application/tag/rpc/internal/svc"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTagLogic {
	return &ListTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTagLogic) ListTag(in *tag.ListTagRequest) (*tag.ListTagReply, error) {
	list, err := l.svcCtx.TagModel.ListTags(l.ctx, in.Cursor, in.Limit+1)
	if err != nil {
		return nil, err
	}

	resp := &tag.ListTagReply{
		Tags:    []*tag.Tag{},
		HasMore: false,
	}

	for i, t := range list {
		if int64(i) < in.Limit {
			resp.Tags = append(resp.Tags, &tag.Tag{
				Id:         int64(t.Id),
				TagName:    t.TagName,
				TagDesc:    t.TagDesc,
				CreateTime: t.CreateTime.Format("2006-01-02 15:04:05"),
			})
		} else {
			resp.HasMore = true
		}
	}
	return resp, nil
}
