package logic

import (
	"context"
	"go_code/zhihu/application/message/rpc/internal/model"
	"time"

	"go_code/zhihu/application/message/rpc/internal/svc"
	"go_code/zhihu/application/message/rpc/types/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *message.SendMessageRequest) (*message.SendMessageReply, error) {
	msg := &model.Message{
		Type:       int64(in.Type),
		BizId:      in.BizId,
		TargetId:   uint64(in.TargetId),
		ReceiverId: uint64(in.ReceiverId),
		Title:      in.Title,
		Content:    in.Content,
		IsRead:     false,
		CreateTime: time.Now(),
	}

	_, err := l.svcCtx.MessageModel.Insert(l.ctx, msg)
	if err != nil {
		return nil, err
	}

	return &message.SendMessageReply{}, nil
}
