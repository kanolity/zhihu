package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/code"
	"go_code/zhihu/application/user/rpc/internal/model"
	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if len(in.Username) == 0 {
		return nil, code.RegisterNameEmpty
	}
	ret, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: in.Username,
		Mobile:   in.Mobile,
		Avatar:   in.Avatar,
		Ctime:    time.Now(),
		Mtime:    time.Now(),
		Password: in.Password,
	})
	if err != nil {
		logx.Errorf("register req:%v err: %v", in, err)
		return nil, err
	}
	userId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId  err: %v", err)
		return nil, err
	}

	return &user.RegisterResponse{UserId: userId}, nil
}
