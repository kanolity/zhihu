package logic

import (
	"context"

	"go_code/zhihu/application/message/rpc/internal/svc"
	"go_code/zhihu/application/message/rpc/types/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *message.GetMessagesRequest) (*message.GetMessagesReply, error) {
	msgs, err := l.svcCtx.MessageModel.ListByReceiver(l.ctx, in.ReceiverId, in.Cursor, in.Limit+1)
	if err != nil {
		return nil, err
	}

	resp := &message.GetMessagesReply{
		Messages: make([]*message.Message, 0, in.Limit),
		HasMore:  false,
	}

	for i, m := range msgs {
		if int64(i) < in.Limit {
			resp.Messages = append(resp.Messages, &message.Message{
				Id:         int64(m.Id),
				Type:       int32(m.Type),
				BizId:      m.BizId,
				TargetId:   int64(m.TargetId),
				Title:      m.Title,
				Content:    m.Content,
				IsRead:     m.IsRead,
				CreateTime: m.CreateTime.Format("2006-01-02 15:04:05"),
			})
		} else {
			resp.HasMore = true
		}
	}
	return resp, nil
}
