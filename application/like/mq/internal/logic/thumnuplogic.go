package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"go_code/zhihu/application/like/mq/internal/model"
	"go_code/zhihu/application/like/mq/internal/svc"
	"go_code/zhihu/application/like/mq/internal/types"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	LikeTypeThumbup       = 0 // 点赞
	LikeTypeDown          = 1 // 点踩
	LikeTypeCancelThumbup = 2 // 取消点赞
	LikeTypeCancelDown    = 3 // 取消点踩
)

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Consume(ctx context.Context, key, val string) error {
	var msg types.ThumbupMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("[Kafka][Thumbup] 解码失败: %v", err)
		return nil
	}

	if msg.Action == "" || msg.ObjId == 0 || msg.UserId == 0 {
		logx.Infof("[Thumbup] 忽略非法消息: %+v", msg)
		return nil
	}

	l.Infof("[Kafka][Thumbup] 收到行为: action=%s biz=%s target=%d user=%d type=%d", msg.Action, msg.BizId, msg.ObjId, msg.UserId, msg.LikeType)

	var err error

	switch msg.Action {
	case "new":
		err = l.svcCtx.LikeModel.Upsert(ctx, &model.Like{
			BizId:    msg.BizId,
			TargetId: uint64(msg.ObjId),
			UserId:   uint64(msg.UserId),
			Type:     int64(msg.LikeType),
			Deleted:  0,
		})
		if err != nil {
			l.Errorf("[Thumbup] Upsert 失败: %v", err)
			return err
		}

		err = l.svcCtx.LikeCountModel.Incr(ctx, msg.BizId, msg.ObjId, msg.LikeType)
		if err != nil {
			l.Errorf("[Thumbup] 计数 Incr 失败: %v", err)
		}

	case "cancel":
		err = l.svcCtx.LikeModel.Upsert(ctx, &model.Like{
			BizId:    msg.BizId,
			TargetId: uint64(msg.ObjId),
			UserId:   uint64(msg.UserId),
			Type:     int64(cancelTargetType(msg.LikeType)), // 存原始点赞类型
			Deleted:  1,
		})
		if err != nil {
			l.Errorf("[Thumbup] Cancel Upsert 失败: %v", err)
			return err
		}

		err = l.svcCtx.LikeCountModel.Decr(ctx, msg.BizId, msg.ObjId, cancelTargetType(msg.LikeType))
		if err != nil {
			l.Errorf("[Thumbup] 计数 Decr 失败: %v", err)
		}

	case "switch":
		// 1. 更新行为表 Type
		err = l.svcCtx.LikeModel.Upsert(ctx, &model.Like{
			BizId:    msg.BizId,
			TargetId: uint64(msg.ObjId),
			UserId:   uint64(msg.UserId),
			Type:     int64(msg.LikeType),
			Deleted:  0,
		})
		if err != nil {
			l.Errorf("[Thumbup] Switch Upsert 失败: %v", err)
			return err
		}

		// 2. 更新点赞计数：增加新类型，减少旧类型
		if msg.OrigType >= 0 {
			_ = l.svcCtx.LikeCountModel.Decr(ctx, msg.BizId, msg.ObjId, msg.OrigType)
		}
		err = l.svcCtx.LikeCountModel.Incr(ctx, msg.BizId, msg.ObjId, msg.LikeType)
		if err != nil {
			l.Errorf("[Thumbup] Switch计数更新失败: %v", err)
		}

	default:
		l.Errorf("[Thumbup] 未知行为 action=%s", msg.Action)
		return nil
	}

	return err
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}

func cancelTargetType(t int32) int32 {
	switch t {
	case LikeTypeCancelThumbup:
		return LikeTypeThumbup
	case LikeTypeCancelDown:
		return LikeTypeDown
	default:
		return -1
	}
}
