package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/types/chat"

	"go_code/zhihu/application/chat/api/internal/svc"
	"go_code/zhihu/application/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq) (resp *types.CreateSessionResp, err error) {
	response, err := l.svcCtx.ChatRpc.CreateSession(l.ctx, &chat.CreateSessionRequest{
		User1Id: req.User1Id,
		User2Id: req.User2Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateSessionResp{
		SessionId: response.SessionId,
	}, nil
}
