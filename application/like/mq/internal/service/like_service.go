package service

import (
	"context"
	"fmt"
	"go_code/zhihu/application/like/mq/internal/model"
	"time"
)

type LikeService struct {
	LikeModel      model.LikeModel
	LikeCountModel model.LikeCountModel
}

// DoThumbup 点赞（写入 DB + 聚合更新）
func (s *LikeService) DoThumbup(ctx context.Context, bizId string, targetId, userId int64, likeType int32) error {
	// 幂等校验
	exist, err := s.LikeModel.Exists(bizId, targetId, userId)
	if err != nil {
		return err
	}
	if exist {
		return nil // 已点赞，跳过
	}

	// 插入点赞明细
	if _, err := s.LikeModel.Insert(ctx, &model.Like{
		BizId:      bizId,
		TargetId:   uint64(targetId),
		UserId:     uint64(userId),
		Type:       int64(likeType),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return fmt.Errorf("插入点赞明细失败: %w", err)
	}

	// 聚合表计数更新
	if err := s.LikeCountModel.Incr(ctx, bizId, targetId, likeType); err != nil {
		return fmt.Errorf("更新聚合计数失败: %w", err)
	}

	return nil
}

func (s *LikeService) CancelThumbup(ctx context.Context, bizId string, targetId, userId int64) error {
	record, err := s.LikeModel.FindByUnique(ctx, bizId, targetId, userId)
	if err != nil {
		return err
	}
	if record == nil {
		return nil
	}

	// 删除点赞记录
	if err := s.LikeModel.Delete(ctx, record.Id); err != nil {
		return fmt.Errorf("删除点赞记录失败: %w", err)
	}

	// 聚合数量减1（）
	err = s.LikeCountModel.Decr(ctx, bizId, targetId, int32(record.Type))

	return nil
}
