package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/internal/model"
	"time"

	"go_code/zhihu/application/chat/rpc/internal/svc"
	"go_code/zhihu/application/chat/rpc/types/chat"

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

func (l *SendMessageLogic) SendMessage(in *chat.SendMessageRequest) (*chat.SendMessageReply, error) {
	msg := &model.ChatMessage{
		SessionId:  uint64(in.SessionId),
		SenderId:   uint64(in.SenderId),
		ReceiverId: uint64(in.ReceiverId),
		Content:    in.Content,
		SendTime:   time.Now(),
	}

	_, err := l.svcCtx.MessageModel.Insert(l.ctx, msg)
	if err != nil {
		return nil, err
	}
	return &chat.SendMessageReply{}, nil
}
