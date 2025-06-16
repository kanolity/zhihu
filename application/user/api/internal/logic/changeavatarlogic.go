package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/user/rpc/userclient"

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
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.UserRpc.ChangeAvatar(l.ctx, &userclient.ChangeAvatarRequest{
		Avatar: req.Avatar,
		UserId: userId,
	})
	return
}
