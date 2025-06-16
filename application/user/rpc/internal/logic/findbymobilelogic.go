package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/pkg/encrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *user.FindByMobileRequest) (*user.FindByMobileResponse, error) {
	mobile, err := encrypt.EncMobile(in.Mobile)
	if err != nil {
		logx.Errorf("encode mobile %s err:%v", in.Mobile, err)
		return nil, err
	}
	user1, err := l.svcCtx.UserModel.FindByMobile(l.ctx, mobile)
	if err != nil {
		logx.Errorf("FindByMobile mobile: %s error: %v", in.Mobile, err)
		return nil, err
	}
	if user1 == nil {
		return &user.FindByMobileResponse{}, nil
	}

	return &user.FindByMobileResponse{
		UserId:   int64(user1.Id),
		Username: user1.Username,
		Avatar:   user1.Avatar,
	}, nil
}
