package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/code"
	"go_code/zhihu/pkg/encrypt"

	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *user.FindByIdRequest) (*user.FindByIdResponse, error) {
	user1, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("FindById userId %s err:%v", in.UserId, err)
		return nil, code.FindUserFailed
	}
	if user1 == nil {
		return &user.FindByIdResponse{}, nil
	}

	user1.Mobile, err = encrypt.DecMobile(user1.Mobile)
	return &user.FindByIdResponse{
		Username: user1.Username,
		UserId:   user1.Id,
		Avatar:   user1.Avatar,
		Mobile:   user1.Mobile,
	}, nil
}
