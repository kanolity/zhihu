package logic

import (
	"context"
	"go_code/zhihu/application/message/rpc/types/message"

	"go_code/zhihu/application/message/api/internal/svc"
	"go_code/zhihu/application/message/api/internal/types"

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
	response, err := l.svcCtx.MessageRpc.GetMessages(l.ctx, &message.GetMessagesRequest{
		ReceiverId: req.ReceiverId,
		Cursor:     req.Cursor,
		Limit:      req.Limit,
	})
	if err != nil {
		return nil, err
	}

	msgs := make([]types.Message, 0, len(response.Messages))
	for _, m := range response.Messages {
		msgs = append(msgs, types.Message{
			Id:         m.Id,
			Type:       m.Type,
			BizId:      m.BizId,
			TargetId:   m.TargetId,
			Title:      m.Title,
			Content:    m.Content,
			IsRead:     m.IsRead,
			CreateTime: m.CreateTime,
		})
	}

	return &types.GetMessagesResp{
		Messages: msgs,
		HasMore:  response.HasMore,
	}, nil
}
