package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/follow/api/internal/svc"
	"go_code/zhihu/application/follow/api/internal/types"
	"go_code/zhihu/application/follow/rpc/followclient"
	"go_code/zhihu/application/user/rpc/types/user"

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
		fanUser, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
			UserId: item.FansUserId})
		if err != nil {
			return nil, err
		}
		newItem := types.FansItem{
			UserId:       item.UserId,
			FansUserId:   item.FansUserId,
			FansUsername: fanUser.Username,
			FansAvatar:   fanUser.Avatar,
			FollowCount:  item.FollowCount,
			FansCount:    item.FansCount,
			CreateTime:   item.CreateTime,
		}
		items = append(items, newItem)
	}
	return &types.FansListResp{
		Items:  items,
		Cursor: response.Cursor,
		IsEnd:  response.IsEnd,
	}, nil
}
