package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/follow/rpc/followclient"

	"go_code/zhihu/application/follow/api/internal/svc"
	"go_code/zhihu/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.FollowListReq) (resp *types.FollowListResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("get user id from context err:%v", err)
		return nil, err
	}
	response, err := l.svcCtx.FollowRpc.FollowList(l.ctx, &followclient.FollowListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		logx.Errorf("call fanslist err:%v", err)
		return nil, err
	}
	items := make([]types.FollowItem, 0, len(response.Items))
	for _, item := range response.Items {
		newItem := types.FollowItem{
			Id:             item.Id,
			FollowedUserId: item.FollowedUserId,
			FansCount:      item.FansCount,
			CreateTime:     item.CreateTime,
		}
		items = append(items, newItem)
	}
	return &types.FollowListResp{
		Items:  items,
		Cursor: response.Cursor,
		IsEnd:  response.IsEnd,
	}, nil
}
