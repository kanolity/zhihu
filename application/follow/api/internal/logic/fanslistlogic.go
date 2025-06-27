package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/follow/api/internal/svc"
	"go_code/zhihu/application/follow/api/internal/types"
	"go_code/zhihu/application/follow/rpc/followclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FansListLogic) FansList(req *types.FansListReq) (resp *types.FansListResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("get user id from context err:%v", err)
		return nil, err
	}
	response, err := l.svcCtx.FollowRpc.FansList(l.ctx, &followclient.FansListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		logx.Errorf("call fanslist err:%v", err)
		return nil, err
	}
	items := make([]types.FansItem, 0, len(response.Items))
	for _, item := range response.Items {
		newItem := types.FansItem{
			UserId:      item.UserId,
			FansUserId:  item.FansUserId,
			FollowCount: item.FollowCount,
			FansCount:   item.FansCount,
			CreateTime:  item.CreateTime,
		}
		items = append(items, newItem)
	}
	return &types.FansListResp{
		Items:  items,
		Cursor: response.Cursor,
		IsEnd:  response.IsEnd,
	}, nil
}
