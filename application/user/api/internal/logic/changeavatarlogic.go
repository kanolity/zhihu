package logic

import (
	"context"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAvatarLogic {
	return &ChangeAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeAvatarLogic) ChangeAvatar(req *types.ChangeAvatarRequest) (resp *types.ChangeAvatarResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
