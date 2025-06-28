package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"go_code/zhihu/application/like/rpc/internal/model"
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
	// 1. 防止重复并发操作，加 Redis 分布式锁
	lockKey := fmt.Sprintf("like:lock:%s:%d:%d", in.BizId, in.ObjId, in.UserId)
	lock := redis.NewRedisLock(l.svcCtx.BizRedis, lockKey)
	lock.SetExpire(3) // 秒
	ok, err := lock.Acquire()
	if err != nil || !ok {
		return nil, errors.New("请求过于频繁")
	}
	defer lock.Release()

	// 2. 查询用户是否已存在点赞行为
	record, err := l.svcCtx.LikeModel.FindByUnique(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		return nil, err
	}

	var action string
	var origType int32 = -1

	switch {
	case record != nil && record.Deleted == 0:
		origType = int32(record.Type)

		switch {
		case in.LikeType == cancelLikeType(origType):
			action = "cancel"
		case in.LikeType != origType:
			action = "switch"
		default:
			cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
			if cnt == nil {
				cnt = &model.LikeCount{}
			}
			return &like.ThumbupResponse{
				BizId:      in.BizId,
				ObjId:      in.ObjId,
				LikeNum:    cnt.LikeNum,
				DislikeNum: cnt.DislikeNum,
			}, nil
		}
	case isCancelType(in.LikeType):
		// 要取消但本身无行为记录，直接返回
		cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
		if cnt == nil {
			cnt = &model.LikeCount{}
		}
		return &like.ThumbupResponse{
			BizId:      in.BizId,
			ObjId:      in.ObjId,
			LikeNum:    cnt.LikeNum,
			DislikeNum: cnt.DislikeNum,
		}, nil
	default:
		action = "new"
	}

	// 3. 构建异步消息结构体
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
		Action:   action,
		OrigType: origType,
	}

	// 4. Kafka 异步投递消息
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal error: %v", err)
			return
		}
		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			l.Logger.Errorf("[Thumbup] kafka push error: %v", err)
		} else {
			l.Logger.Infof("[Thumbup] action=%s user=%d obj=%d type=%d → 推送成功", action, in.UserId, in.ObjId, in.LikeType)
		}
	})

	// 5. 返回当前计数（可能是旧值）
	cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
	if cnt == nil {
		cnt = &model.LikeCount{}
	}
	return &like.ThumbupResponse{
		BizId:      in.BizId,
		ObjId:      in.ObjId,
		LikeNum:    cnt.LikeNum,
		DislikeNum: cnt.DislikeNum,
	}, nil
}

func isCancelType(t int32) bool {
	return t == 2 || t == 3 // CancelThumbup / CancelDown
}

func cancelLikeType(orig int32) int32 {
	switch orig {
	case 0:
		return 2
	case 1:
		return 3
	default:
		return -1
	}
}
