package logic

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	existingUser.Avatar = in.Avatar
	existingUser.Mtime = time.Now()

	err = l.svcCtx.UserModel.Update(l.ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("头像更新失败: %w", err)
	}

	return &user.ChangeAvatarResponse{}, nil
}
