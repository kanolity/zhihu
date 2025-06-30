package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go_code/zhihu/application/user/rpc/internal/code"
	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/pkg/encrypt"
)

var cacheBeyondUserUserIdPrefix = "cache:beyondUser:user:id:"

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
	if in.OldPassword != (*user1).Password {
		fmt.Println("in.OldPassword=", in.OldPassword)
		fmt.Println("user1.Password", (*user1).Password)
		fmt.Printf("%#v\n", user1)
		return nil, code.ChangePasswordWrong
	}
	user1.Password = encrypt.EncPassword(in.NewPassword)
	err = l.svcCtx.UserModel.Update(l.ctx, user1)
	if err != nil {
		logx.Errorf("ChangePassword userId:%v err: %v", user1.Id, err)
		return nil, code.ChangePasswordFailed
	}

	//删除缓存
	cacheKey := fmt.Sprintf("%s%v", cacheBeyondUserUserIdPrefix, user1.Id)
	_, err = l.svcCtx.BizRedis.DelCtx(l.ctx, cacheKey)
	if err != nil {
		logx.Errorf("Delete user cache failed, key=%s: %v", cacheKey, err)
	}

	return &user.ChangePasswordResponse{}, nil
}
