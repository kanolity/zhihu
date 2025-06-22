package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeCountModel = (*customLikeCountModel)(nil)

type (
	// LikeCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeCountModel.
	LikeCountModel interface {
		likeCountModel
	}

	customLikeCountModel struct {
		*defaultLikeCountModel
	}
)

// NewLikeCountModel returns a model for the database table.
func NewLikeCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeCountModel {
	return &customLikeCountModel{
		defaultLikeCountModel: newLikeCountModel(conn, c, opts...),
	}
}

func (m *defaultLikeCountModel) Incr(ctx context.Context, bizId string, targetId int64, likeType int32) error {
	field := "like_num"
	if likeType == 1 {
		field = "dislike_num"
	}
	query := fmt.Sprintf(`
         INSERT INTO %s (biz_id, target_id, %s)
        VALUES (?, ?, 1)
        ON DUPLICATE KEY UPDATE
            %s = %s + 1,
            update_time = NOW()
    `, m.table, field, field, field)
	_, err := m.CachedConn.ExecNoCacheCtx(ctx, query, bizId, targetId)
	return err
}

func (m *defaultLikeCountModel) Decr(ctx context.Context, bizId string, targetId int64, likeType int32) error {
	field := "like_num"
	if likeType == 4 {
		field = "dislike_num"
	}
	query := fmt.Sprintf(`
        INSERT INTO %s (biz_id, target_id, %s)
        VALUES (?, ?, 0)
        ON DUPLICATE KEY UPDATE
            %s = GREATEST(%s - 1, 0),
            update_time = NOW()
    `, m.table, field, field, field)
	_, err := m.CachedConn.ExecNoCacheCtx(ctx, query, bizId, targetId)
	return err
}
