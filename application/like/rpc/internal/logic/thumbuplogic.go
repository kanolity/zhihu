package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"go_code/zhihu/application/like/rpc/internal/svc"
	"go_code/zhihu/application/like/rpc/types"
	"go_code/zhihu/application/like/rpc/types/like"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *like.ThumbupRequest) (*like.ThumbupResponse, error) {
	// 1. 加锁防止并发重复点赞
	lockKey := fmt.Sprintf("like:lock:%s:%d:%d", in.BizId, in.ObjId, in.UserId)
	lock := redis.NewRedisLock(l.svcCtx.BizRedis, lockKey)
	lock.SetExpire(3) // 秒级别
	ok, err := lock.Acquire()
	if err != nil || !ok {
		return nil, errors.New("请求过于频繁")
	}
	defer lock.Release()

	// 2. 判断用户是否已经点过（从缓存/DB判断）
	exist, err := l.svcCtx.LikeModel.Exists(in.BizId, in.ObjId, in.UserId)
	if err != nil {
		return nil, err
	}
	if exist {
		cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
		return &like.ThumbupResponse{
			BizId:      in.BizId,
			ObjId:      in.ObjId,
			LikeNum:    cnt.LikeNum,
			DislikeNum: cnt.DislikeNum,
		}, nil
	}

	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}
	// 3. 异步投递点赞行为到 Kafka
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Logger.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
		}
	})

	// 4. 返回点赞后的最新数量
	cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
	return &like.ThumbupResponse{
		BizId:      in.BizId,
		ObjId:      in.ObjId,
		LikeNum:    cnt.LikeNum,
		DislikeNum: cnt.DislikeNum,
	}, nil
}
