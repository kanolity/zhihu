package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/follow/rpc/followclient"

	"go_code/zhihu/application/follow/api/internal/svc"
	"go_code/zhihu/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfollowLogic) Unfollow(req *types.UnFollowReq) (resp *types.UnFollowResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("get user id from context err:%v", err)
		return nil, err
	}

	_, err = l.svcCtx.FollowRpc.UnFollow(l.ctx, &followclient.UnFollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		logx.Errorf("unfollow err:%v", err)
		return nil, err
	}

	return &types.UnFollowResp{}, nil
}
