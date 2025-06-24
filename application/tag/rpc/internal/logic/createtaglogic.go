package logic

import (
	"context"
	"go_code/zhihu/application/tag/rpc/internal/model"
	"time"

	"go_code/zhihu/application/tag/rpc/internal/svc"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTagLogic) CreateTag(in *tag.CreateTagRequest) (*tag.CreateTagReply, error) {
	now := time.Now()

	t := &model.Tag{
		TagName:    in.TagName,
		TagDesc:    in.TagDesc,
		CreateTime: now,
		UpdateTime: now,
	}

	res, err := l.svcCtx.TagModel.Insert(l.ctx, t)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &tag.CreateTagReply{Id: id}, nil
}
