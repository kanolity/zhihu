package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/types/chat"

	"go_code/zhihu/application/chat/api/internal/svc"
	"go_code/zhihu/application/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageReq) (resp *types.SendMessageResp, err error) {
	_, err = l.svcCtx.ChatRpc.SendMessage(l.ctx, &chat.SendMessageRequest{
		SessionId:  req.SessionId,
		SenderId:   req.SenderId,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
	})
	return &types.SendMessageResp{}, err
}
