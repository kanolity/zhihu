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
	// todo: add your logic here and delete this line

	return &user.FindByUsernameResponse{}, nil
}
