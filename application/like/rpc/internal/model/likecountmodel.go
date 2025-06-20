package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeCountModel = (*customLikeCountModel)(nil)
var cacheLikeCountBizTargetIdPrefix = "likecount:biz:%s:target:%d"

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
        INSERT INTO %s (biz_id, target_id, %s) VALUES (?, ?, 1)
        ON DUPLICATE KEY UPDATE %s = %s + 1
    `, m.table, field, field, field)
	_, err := m.CachedConn.ExecNoCacheCtx(ctx, query, bizId, targetId)
	return err
}

func (m *defaultLikeCountModel) FindByBizTarget(ctx context.Context, bizId string, targetId int64) (*LikeCount, error) {
	key := fmt.Sprintf(cacheLikeCountBizTargetIdPrefix, bizId, targetId)
	var resp LikeCount

	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("SELECT %s FROM %s WHERE biz_id = ? AND target_id = ?", likeCountRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, bizId, targetId)
	})
	if err == sqlc.ErrNotFound {
		return nil, nil
	}
	return &resp, err
}
