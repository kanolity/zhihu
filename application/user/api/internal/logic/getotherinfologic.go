package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/types/user"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOtherInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOtherInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOtherInfoLogic {
	return &GetOtherInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOtherInfoLogic) GetOtherInfo(req *types.GetOtherInfoRequest) (resp *types.GetOtherInfoResponse, err error) {
	if req.UserId == 0 {
		return &types.GetOtherInfoResponse{}, nil
	}
	u, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.Errorf("FindById userId: %d error: %v", req.UserId, err)
		return nil, err
	}

	return &types.GetOtherInfoResponse{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, nil
}
