package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/internal/svc"
	"go_code/zhihu/application/chat/rpc/types/chat"

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

func (l *GetMessagesLogic) GetMessages(in *chat.GetMessagesRequest) (*chat.GetMessagesReply, error) {
	list, err := l.svcCtx.MessageModel.GetMessagesBySession(l.ctx, in.SessionId, in.Cursor, in.Limit+1)
	if err != nil {
		return nil, err
	}

	resp := &chat.GetMessagesReply{
		Messages: make([]*chat.Message, 0, in.Limit),
		HasMore:  false,
	}

	var unreadIDs []int64

	for i, msg := range list {
		if int64(i) < in.Limit {
			resp.Messages = append(resp.Messages, &chat.Message{
				Id:         int64(msg.Id),
				SenderId:   int64(msg.SenderId),
				ReceiverId: int64(msg.ReceiverId),
				Content:    msg.Content,
				IsRead:     msg.IsRead,
				SendTime:   msg.SendTime.Format("2006-01-02 15:04:05"),
			})
			if msg.ReceiverId == uint64(in.UserId) && !msg.IsRead {
				unreadIDs = append(unreadIDs, int64(msg.Id))
			}
		} else {
			resp.HasMore = true
		}
	}
	//标记为已读
	if len(unreadIDs) > 0 {
		if err := l.svcCtx.MessageModel.MarkAsReadBatch(l.ctx, unreadIDs); err != nil {
			logx.Errorf("mark messages read error: %v", err)
		}
	}
	return resp, nil
}
