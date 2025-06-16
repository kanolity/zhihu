package logic

import (
	"context"
	"encoding/json"

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
	userId, err := l.ctx.Value("userId").(json.Number).Int64()

	return &user.ChangeAvatarResponse{}, nil
}
