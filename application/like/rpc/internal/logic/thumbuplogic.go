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
	// 1. 加锁防止并发重复操作
	lockKey := fmt.Sprintf("like:lock:%s:%d:%d", in.BizId, in.ObjId, in.UserId)
	lock := redis.NewRedisLock(l.svcCtx.BizRedis, lockKey)
	lock.SetExpire(3) // 秒
	ok, err := lock.Acquire()
	if err != nil || !ok {
		return nil, errors.New("请求过于频繁")
	}
	defer lock.Release()

	// 2. 获取用户现有行为记录
	record, err := l.svcCtx.LikeModel.FindByUnique(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		return nil, err
	}

	// 3. 状态决策逻辑
	if record != nil && record.Deleted == 0 {
		if in.LikeType == cancelLikeType(int32(record.Type)) {
			// 👇 用户希望“取消”之前的点赞/点踩
			// 异步投递取消行为
		} else if int32(record.Type) != in.LikeType {
			// 👇 用户切换了行为（点踩 → 点赞）
			// 异步投递“更新行为类型”
		} else {
			// 👇 幂等重复行为，直接返回现有状态
			cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
			return &like.ThumbupResponse{
				BizId:      in.BizId,
				ObjId:      in.ObjId,
				LikeNum:    cnt.LikeNum,
				DislikeNum: cnt.DislikeNum,
			}, nil
		}
	} else {
		if isCancelType(in.LikeType) {
			// 👇 用户尝试取消不存在的行为 → 无需处理，直接返回
			cnt, _ := l.svcCtx.LikeCountModel.FindByBizTarget(l.ctx, in.BizId, in.ObjId)
			return &like.ThumbupResponse{
				BizId:      in.BizId,
				ObjId:      in.ObjId,
				LikeNum:    cnt.LikeNum,
				DislikeNum: cnt.DislikeNum,
			}, nil
		}
		// 👇 首次点赞/点踩行为，继续投递
	}

	// 4. 构造异步消息
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal error: %v", err)
			return
		}
		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			l.Logger.Errorf("[Thumbup] push error: %v", err)
		}
	})

	// 5. 返回当前数据（可能是旧值）
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
