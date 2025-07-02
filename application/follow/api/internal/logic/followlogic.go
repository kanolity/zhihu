package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"go_code/zhihu/application/follow/rpc/followclient"

	"go_code/zhihu/application/follow/api/internal/svc"
	"go_code/zhihu/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowReq) (resp *types.FollowResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	fmt.Printf("%#v\n", req)
	fmt.Printf("userid=%d\n", userId)
	if err != nil {
		logx.Errorf("get user id from context err:%v", err)
		return nil, err
	}

	_, err = l.svcCtx.FollowRpc.Follow(l.ctx, &followclient.FollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		logx.Errorf("follow err:%v", err)
		return nil, err
	}

	return &types.FollowResp{}, nil
}
