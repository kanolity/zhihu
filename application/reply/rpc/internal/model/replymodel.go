package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReplyModel = (*customReplyModel)(nil)

type (
	// ReplyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReplyModel.
	ReplyModel interface {
		replyModel
		withSession(session sqlx.Session) ReplyModel
		ListByTarget(ctx context.Context, bizId string, targetId, cursor, limit int64) ([]*Reply, error)
	}

	customReplyModel struct {
		*defaultReplyModel
	}
)

// NewReplyModel returns a model for the database table.
func NewReplyModel(conn sqlx.SqlConn) ReplyModel {
	return &customReplyModel{
		defaultReplyModel: newReplyModel(conn),
	}
}

func (m *customReplyModel) withSession(session sqlx.Session) ReplyModel {
	return NewReplyModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultReplyModel) ListByTarget(ctx context.Context, bizId string, targetId, cursor, limit int64) ([]*Reply, error) {
	var args []interface{}
	var query string

	if cursor > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE biz_id = ? AND target_id = ? AND id < ? AND status = 0 ORDER BY id DESC LIMIT ?", replyRows, m.table)
		args = []interface{}{bizId, targetId, cursor, limit}
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE biz_id = ? AND target_id = ? AND status = 0 ORDER BY id DESC LIMIT ?", replyRows, m.table)
		args = []interface{}{bizId, targetId, limit}
	}

	var replies []*Reply
	if err := m.conn.QueryRowsCtx(ctx, &replies, query, args...); err != nil {
		return nil, err
	}
	return replies, nil
}
