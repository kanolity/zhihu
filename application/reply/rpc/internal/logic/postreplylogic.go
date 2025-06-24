package logic

import (
	"context"
	"go_code/zhihu/application/reply/rpc/internal/model"
	"time"

	"go_code/zhihu/application/reply/rpc/internal/svc"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostReplyLogic {
	return &PostReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostReplyLogic) PostReply(in *reply.PostReplyRequest) (*reply.PostReplyReply, error) {
	now := time.Now()

	r := &model.Reply{
		BizId:         in.BizId,
		TargetId:      uint64(in.TargetId),
		ReplyUserId:   uint64(in.ReplyUserId),
		BeReplyUserId: uint64(in.BeReplyUserId),
		ParentId:      uint64(in.ParentId),
		Content:       in.Content,
		Status:        0,
		LikeNum:       0,
		CreateTime:    now,
		UpdateTime:    now,
	}

	ret, err := l.svcCtx.ReplyModel.Insert(l.ctx, r)
	if err != nil {
		return nil, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId  err: %v", err)
		return nil, err
	}

	// 可选：更新 reply_count 表
	_ = l.svcCtx.ReplyCountModel.IncreaseCount(l.ctx, in.BizId, in.TargetId, in.ParentId == 0)

	return &reply.PostReplyReply{Id: id}, nil
}
