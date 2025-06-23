package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/types/chat"

	"go_code/zhihu/application/chat/api/internal/svc"
	"go_code/zhihu/application/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessagesLogic) GetMessages(req *types.GetMessagesReq) (resp *types.GetMessagesResp, err error) {
	response, err := l.svcCtx.ChatRpc.GetMessages(l.ctx, &chat.GetMessagesRequest{
		SessionId: req.SessionId,
		Cursor:    req.Cursor,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]types.Message, 0, len(response.Messages))
	for _, m := range response.Messages {
		messages = append(messages, types.Message{
			Id:         m.Id,
			SenderId:   m.SenderId,
			ReceiverId: m.ReceiverId,
			Content:    m.Content,
			IsRead:     m.IsRead,
			SendTime:   m.SendTime,
		})
	}

	return &types.GetMessagesResp{
		Messages: messages,
		HasMore:  response.HasMore,
	}, nil
}
