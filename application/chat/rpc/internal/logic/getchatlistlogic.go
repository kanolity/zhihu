package logic

import (
	"context"

	"go_code/zhihu/application/chat/rpc/internal/svc"
	"go_code/zhihu/application/chat/rpc/types/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatListLogic {
	return &GetChatListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatListLogic) GetChatList(in *chat.GetChatListRequest) (*chat.GetChatListResponse, error) {
	sessions, err := l.svcCtx.SessionModel.FindUserSessions(l.ctx, in.UserId, in.Cursor, int(in.Limit+1))
	if err != nil {
		return nil, err
	}

	var resp chat.GetChatListResponse
	if len(sessions) > int(in.Limit) {
		resp.HasMore = true
		sessions = sessions[:in.Limit]
	}

	for _, sess := range sessions {
		// 查询该会话的最新消息
		msg, err := l.svcCtx.MessageModel.FindLatestMessageBySession(l.ctx, int64(sess.Id))
		if err != nil {
			logx.Errorf("failed to get latest message for session %d: %v", sess.Id, err)
			continue
		}

		// 获取对方 ID（双人会话中 userId ≠ self）
		var targetId int64
		if in.UserId == int64(sess.User1Id) {
			targetId = int64(sess.User2Id)
		} else {
			targetId = int64(sess.User1Id)
		}

		resp.ChatList = append(resp.ChatList, &chat.ChatList{
			TargetUserId: targetId,
			LatestMessage: &chat.Message{
				Id:         int64(msg.Id),
				SenderId:   int64(msg.SenderId),
				ReceiverId: int64(msg.ReceiverId),
				Content:    msg.Content,
				SendTime:   msg.SendTime.Format("2006-01-02 15:04:05"),
				IsRead:     msg.IsRead,
			},
		})
	}

	return &resp, nil
}
