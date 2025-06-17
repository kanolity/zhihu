package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/code"
	"time"

	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAvatarLogic {
	return &ChangeAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeAvatarLogic) ChangeAvatar(in *user.ChangeAvatarRequest) (*user.ChangeAvatarResponse, error) {
	existingUser, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		return nil, code.FindUserFailed
	}

	existingUser.Avatar = in.Avatar
	existingUser.Mtime = time.Now()

	err = l.svcCtx.UserModel.Update(l.ctx, existingUser)
	if err != nil {
		logx.Errorf("ChangeAvatar userId:%v err:%v", existingUser.Id, err)
		return nil, code.ChangeAvatarFailed
	}

	return &user.ChangeAvatarResponse{}, nil
}
