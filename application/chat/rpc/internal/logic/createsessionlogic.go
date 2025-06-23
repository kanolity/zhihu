package logic

import (
	"context"
	"go_code/zhihu/application/chat/rpc/internal/model"

	"go_code/zhihu/application/chat/rpc/internal/svc"
	"go_code/zhihu/application/chat/rpc/types/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *chat.CreateSessionRequest) (*chat.CreateSessionReply, error) {
	// 避免重复对话
	exist, err := l.svcCtx.SessionModel.FindByUserPair(l.ctx, in.User1Id, in.User2Id)
	if err == nil && exist != nil {
		return &chat.CreateSessionReply{SessionId: int64(exist.Id)}, nil
	}

	ret, err := l.svcCtx.SessionModel.Insert(l.ctx, &model.ChatSession{
		User1Id: uint64(in.User1Id),
		User2Id: uint64(in.User2Id),
	})
	if err != nil {
		return nil, err
	}
	newId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId  err: %v", err)
		return nil, err
	}
	return &chat.CreateSessionReply{SessionId: newId}, nil
}
