package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		withSession(session sqlx.Session) MessageModel
		ListByReceiver(ctx context.Context, receiverId, cursor, limit int64) ([]*Message, error)
		MarkRead(ctx context.Context, id int64) error
		ListUnread(ctx context.Context, receiverId, cursor, limit int64) ([]*Message, error)
		CountUnread(ctx context.Context, receiverId int64) (int64, error)
		MarkAllRead(ctx context.Context, receiverId int64) error
	}

	customMessageModel struct {
		*defaultMessageModel
	}
)

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn sqlx.SqlConn) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn),
	}
}

func (m *customMessageModel) withSession(session sqlx.Session) MessageModel {
	return NewMessageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMessageModel) ListByReceiver(ctx context.Context, receiverId, cursor, limit int64) ([]*Message, error) {
	var args []interface{}
	var query string

	if cursor > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE receiver_id = ? AND id < ? ORDER BY id DESC LIMIT ?", messageRows, m.table)
		args = []interface{}{receiverId, cursor, limit}
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE receiver_id = ? ORDER BY id DESC LIMIT ?", messageRows, m.table)
		args = []interface{}{receiverId, limit}
	}

	var list []*Message
	if err := m.conn.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *defaultMessageModel) MarkRead(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET is_read = 1 WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMessageModel) ListUnread(ctx context.Context, receiverId, cursor, limit int64) ([]*Message, error) {
	var args []interface{}
	var query string

	if cursor > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE receiver_id = ? AND is_read = 0 AND id < ? ORDER BY id DESC LIMIT ?", messageRows, m.table)
		args = []interface{}{receiverId, cursor, limit}
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE receiver_id = ? AND is_read = 0 ORDER BY id DESC LIMIT ?", messageRows, m.table)
		args = []interface{}{receiverId, limit}
	}

	var list []*Message
	err := m.conn.QueryRowsCtx(ctx, &list, query, args...)
	return list, err
}

func (m *defaultMessageModel) CountUnread(ctx context.Context, receiverId int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE receiver_id = ? AND is_read = 0", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, receiverId)
	return count, err
}

func (m *defaultMessageModel) MarkAllRead(ctx context.Context, receiverId int64) error {
	query := fmt.Sprintf("UPDATE %s SET is_read = 1 WHERE receiver_id = ? AND is_read = 0", m.table)
	_, err := m.conn.ExecCtx(ctx, query, receiverId)
	return err
}
