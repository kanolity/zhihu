package logic

import (
	"context"
	"go_code/zhihu/application/follow/rpc/internal/code"
	"go_code/zhihu/application/follow/rpc/internal/model"
	"go_code/zhihu/application/follow/rpc/types"
	"gorm.io/gorm"
	"strconv"

	"go_code/zhihu/application/follow/rpc/internal/svc"
	"go_code/zhihu/application/follow/rpc/types/follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UnFollow 取消关注
func (l *UnFollowLogic) UnFollow(in *follow.UnFollowRequest) (*follow.UnFollowResponse, error) {
	if in.UserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowedUserIdEmpty
	}

	follow1, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[UnFollow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	if follow1 == nil {
		return &follow.UnFollowResponse{}, nil
	}
	if follow1.FollowStatus == types.FollowStatusUnfollow {
		return &follow.UnFollowResponse{}, nil
	}

	// 事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err := model.NewFollowModel(tx).UpdateFields(l.ctx, follow1.ID, map[string]interface{}{
			"follow_status": types.FollowStatusUnfollow,
		})
		if err != nil {
			return err
		}
		err = model.NewFollowCountModel(tx).DecrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).DecrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[UnFollow] Transaction error: %v", err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFollowKey(in.UserId), strconv.FormatInt(in.FollowedUserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] BizRedis.ZremCtx error: %v", err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFansKey(in.FollowedUserId), strconv.FormatInt(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] BizRedis.ZremCtx error: %v", err)
		return nil, err
	}

	return &follow.UnFollowResponse{}, nil
}
