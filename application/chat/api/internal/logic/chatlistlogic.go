package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/chat/rpc/types/chat"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/application/user/rpc/userclient"

	"go_code/zhihu/application/chat/api/internal/svc"
	"go_code/zhihu/application/chat/api/internal/types"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatListLogic) ChatList(req *types.GetChatListReq) (resp *types.GetChatListResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	// 调用 chatRpc 获取会话列表
	rpcResp, err := l.svcCtx.ChatRpc.GetChatList(l.ctx, &chat.GetChatListRequest{
		UserId: userId,
		Cursor: req.Cursor,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	// 抽取对方 ID 集合，准备调用 userRpc
	targetIds := lo.Map(rpcResp.ChatList, func(c *chat.ChatList, _ int) int64 {
		return c.TargetUserId
	})

	userInfoResp, err := l.svcCtx.UserRpc.BatchGetUsers(l.ctx, &userclient.BatchGetUsersRequest{
		UserIds: targetIds,
	})
	if err != nil {
		return nil, err
	}

	// 建立 userId → 用户名/头像 映射表
	userMap := make(map[int64]*user.UserInfo)
	for _, u := range userInfoResp.Users {
		userMap[u.Id] = u
	}

	// 拼接返回结构
	chats := lo.Map(rpcResp.ChatList, func(item *chat.ChatList, _ int) types.ChatList {
		u := userMap[item.TargetUserId]
		return types.ChatList{
			TargetUserId:  item.TargetUserId,
			Username:      u.GetUsername(),
			Avatar:        u.GetAvatar(),
			LatestMessage: convertMessage(item.LatestMessage),
		}
	})

	return &types.GetChatListResponse{
		Chats:   chats,
		HasMore: rpcResp.HasMore,
	}, nil
}

func convertMessage(m *chat.Message) types.Message {
	if m == nil {
		return types.Message{}
	}
	return types.Message{
		Id:         m.Id,
		SenderId:   m.SenderId,
		ReceiverId: m.ReceiverId,
		Content:    m.Content,
		SendTime:   m.SendTime,
		IsRead:     m.IsRead,
	}
}
