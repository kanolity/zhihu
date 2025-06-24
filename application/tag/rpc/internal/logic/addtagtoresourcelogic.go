package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"go_code/zhihu/application/tag/rpc/internal/model"
	"time"

	"go_code/zhihu/application/tag/rpc/internal/svc"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTagToResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTagToResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagToResourceLogic {
	return &AddTagToResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddTagToResourceLogic) AddTagToResource(in *tag.AddTagToResourceRequest) (*tag.AddTagToResourceReply, error) {
	now := time.Now()

	tr := &model.TagResource{
		BizId:      in.BizId,
		TargetId:   uint64(in.TargetId),
		TagId:      uint64(in.TagId),
		UserId:     uint64(in.UserId),
		CreateTime: now,
		UpdateTime: now,
	}

	// 可选：幂等判断，避免重复绑定
	exist, err := l.svcCtx.TagResourceModel.FindOneByUniqueKey(l.ctx, in.BizId, in.TargetId, in.TagId)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, err
	}
	if exist != nil {
		return &tag.AddTagToResourceReply{}, nil
	}

	_, err = l.svcCtx.TagResourceModel.Insert(l.ctx, tr)
	return &tag.AddTagToResourceReply{}, err
}
