package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TagResourceModel = (*customTagResourceModel)(nil)

type (
	// TagResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagResourceModel.
	TagResourceModel interface {
		tagResourceModel
		withSession(session sqlx.Session) TagResourceModel
		FindByBizTarget(ctx context.Context, bizId string, targetId int64) ([]*TagResource, error)
		FindOneByUniqueKey(ctx context.Context, bizId string, targetId, tagId int64) (*TagResource, error)
	}

	customTagResourceModel struct {
		*defaultTagResourceModel
	}
)

// NewTagResourceModel returns a model for the database table.
func NewTagResourceModel(conn sqlx.SqlConn) TagResourceModel {
	return &customTagResourceModel{
		defaultTagResourceModel: newTagResourceModel(conn),
	}
}

func (m *customTagResourceModel) withSession(session sqlx.Session) TagResourceModel {
	return NewTagResourceModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultTagResourceModel) FindByBizTarget(ctx context.Context, bizId string, targetId int64) ([]*TagResource, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE biz_id = ? AND target_id = ? ORDER BY id ASC`, tagResourceRows, m.table)
	var res []*TagResource
	err := m.conn.QueryRowsCtx(ctx, &res, query, bizId, targetId)
	return res, err
}

func (m *defaultTagResourceModel) FindOneByUniqueKey(ctx context.Context, bizId string, targetId, tagId int64) (*TagResource, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE biz_id = ? AND target_id = ? AND tag_id = ? LIMIT 1`, tagResourceRows, m.table)
	var tr TagResource
	err := m.conn.QueryRowCtx(ctx, &tr, query, bizId, targetId, tagId)
	if err == sqlc.ErrNotFound {
		return nil, nil
	}
	return &tr, err
}
