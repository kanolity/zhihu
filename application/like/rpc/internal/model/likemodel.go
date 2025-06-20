package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeModel = (*customLikeModel)(nil)

type (
	// LikeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeModel.
	LikeModel interface {
		likeModel
	}

	customLikeModel struct {
		*defaultLikeModel
	}
)

// NewLikeModel returns a model for the database table.
func NewLikeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeModel {
	return &customLikeModel{
		defaultLikeModel: newLikeModel(conn, c, opts...),
	}
}

func (m *defaultLikeModel) Exists(bizId string, targetId, userId int64) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE biz_id = ? AND target_id = ? AND user_id = ?", m.table)
	var count int64
	err := m.CachedConn.QueryRowNoCache(&count, query, bizId, targetId, userId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (m *customLikeModel) FindByUnique(ctx context.Context, bizId string, targetId, userId int64) (*Like, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE biz_id = ? AND target_id = ? AND user_id = ? LIMIT 1", likeRows, m.table)
	var resp Like
	err := m.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, bizId, targetId, userId)
	if err == sqlc.ErrNotFound {
		return nil, nil
	}
	return &resp, err
}
