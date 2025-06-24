package logic

import (
	"context"

	"go_code/zhihu/application/message/rpc/internal/svc"
	"go_code/zhihu/application/message/rpc/types/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAsReadLogic {
	return &MarkAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkAsReadLogic) MarkAsRead(in *message.MarkAsReadRequest) (*message.MarkAsReadReply, error) {
	err := l.svcCtx.MessageModel.MarkRead(l.ctx, in.Id)
	return &message.MarkAsReadReply{}, err
}
