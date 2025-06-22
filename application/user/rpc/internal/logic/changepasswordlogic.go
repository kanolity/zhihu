package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/code"
	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/pkg/encrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	user1, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, code.FindUserFailed
	}
	in.OldPassword = encrypt.EncPassword(in.OldPassword)
	if in.OldPassword != user1.Password {
		return nil, code.ChangePasswordWrong
	}
	user1.Password = encrypt.EncPassword(in.NewPassword)
	err = l.svcCtx.UserModel.Update(l.ctx, user1)
	if err != nil {
		logx.Errorf("ChangePassword userId:%v err: %v", user1.Id, err)
		return nil, code.ChangePasswordFailed
	}

	return &user.ChangePasswordResponse{}, nil
}
