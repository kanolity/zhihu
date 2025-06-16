package logic

import (
	"context"

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
	user1, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		logx.Errorf("FindById userId %s err:%v", in.UserId, err)
		return nil, err
	}
	if user1 == nil {
		return &user.FindByIdResponse{}, nil
	}

	return &user.FindByIdResponse{
		Username: user1.Username,
		UserId:   int64(user1.Id),
		Avatar:   user1.Avatar,
	}, nil
}
