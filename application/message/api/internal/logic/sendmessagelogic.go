package logic

import (
	"context"
	"go_code/zhihu/application/message/rpc/types/message"

	"go_code/zhihu/application/message/api/internal/svc"
	"go_code/zhihu/application/message/api/internal/types"

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
	_, err = l.svcCtx.MessageRpc.SendMessage(l.ctx, &message.SendMessageRequest{
		Type:       req.Type,
		BizId:      req.BizId,
		TargetId:   req.TargetId,
		ReceiverId: req.ReceiverId,
		Title:      req.Title,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &types.SendMessageResp{}, nil
}
