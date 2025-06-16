package logic

import (
	"context"

	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByUsernameLogic {
	return &FindByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByUsernameLogic) FindByUsername(in *user.FindByUsernameRequest) (*user.FindByUsernameResponse, error) {
	user1, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		logx.Errorf("FindByUsername username %s err:%v", in.Username, err)
		return nil, err
	}
	if user1 == nil {
		return &user.FindByUsernameResponse{}, nil
	}

	return &user.FindByUsernameResponse{
		Username: user1.Username,
		UserId:   int64(user1.Id),
		Avatar:   user1.Avatar,
	}, nil
}
