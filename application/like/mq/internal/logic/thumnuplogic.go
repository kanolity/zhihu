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

	l.Infof("[Kafka][Thumbup] 收到 → biz=%s target=%d user=%d type=%d", msg.BizId, msg.ObjId, msg.UserId, msg.LikeType)

	// 点赞行为记录处理
	var likeTypeToWrite int64
	var deleted int64 = 0

	switch msg.LikeType {
	case LikeTypeThumbup, LikeTypeDown:
		likeTypeToWrite = int64(msg.LikeType)
		deleted = 0
	case LikeTypeCancelThumbup:
		likeTypeToWrite = int64(LikeTypeThumbup)
		deleted = 1
	case LikeTypeCancelDown:
		likeTypeToWrite = int64(LikeTypeDown)
		deleted = 1
	default:
		l.Errorf("[Kafka][Thumbup] 无效 LikeType: %d", msg.LikeType)
		return nil
	}

	// Upsert 行为表
	err := l.svcCtx.LikeModel.Upsert(ctx, &model.Like{
		BizId:    msg.BizId,
		TargetId: uint64(msg.ObjId),
		UserId:   uint64(msg.UserId),
		Type:     likeTypeToWrite,
		Deleted:  deleted,
	})
	if err != nil {
		l.Errorf("[Thumbup] Upsert 行为失败: %v", err)
		return err
	}

	// 点赞计数更新
	switch msg.LikeType {
	case LikeTypeThumbup, LikeTypeDown:
		err = l.svcCtx.LikeCountModel.Incr(ctx, msg.BizId, msg.ObjId, msg.LikeType)
	case LikeTypeCancelThumbup:
		err = l.svcCtx.LikeCountModel.Decr(ctx, msg.BizId, msg.ObjId, LikeTypeThumbup)
	case LikeTypeCancelDown:
		err = l.svcCtx.LikeCountModel.Decr(ctx, msg.BizId, msg.ObjId, LikeTypeDown)
	}

	if err != nil {
		l.Errorf("[Thumbup] 点赞计数更新失败: %v", err)
	}

	return err
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}
