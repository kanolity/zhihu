package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ ChatMessageModel = (*customChatMessageModel)(nil)

type (
	// ChatMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatMessageModel.
	ChatMessageModel interface {
		chatMessageModel
		withSession(session sqlx.Session) ChatMessageModel
		GetMessagesBySession(ctx context.Context, sessionId, cursor, limit int64) ([]*ChatMessage, error)
		MarkAsRead(ctx context.Context, messageId int64) error
		MarkAsReadBatch(ctx context.Context, ids []int64) error
		FindLatestMessageBySession(ctx context.Context, sessionId int64) (*ChatMessage, error)
	}

	customChatMessageModel struct {
		*defaultChatMessageModel
	}
)

// NewChatMessageModel returns a model for the database table.
func NewChatMessageModel(conn sqlx.SqlConn) ChatMessageModel {
	return &customChatMessageModel{
		defaultChatMessageModel: newChatMessageModel(conn),
	}
}

func (m *customChatMessageModel) withSession(session sqlx.Session) ChatMessageModel {
	return NewChatMessageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultChatMessageModel) GetMessagesBySession(ctx context.Context, sessionId, cursor, limit int64) ([]*ChatMessage, error) {
	var (
		query string
		args  []interface{}
	)
	if cursor > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE session_id = ? AND id < ? ORDER BY id DESC LIMIT ?", chatMessageRows, m.table)
		args = []interface{}{sessionId, cursor, limit}
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE session_id = ? ORDER BY id DESC LIMIT ?", chatMessageRows, m.table)
		args = []interface{}{sessionId, limit}
	}

	var messages []*ChatMessage
	err := m.conn.QueryRowsCtx(ctx, &messages, query, args...)
	if err == sqlc.ErrNotFound {
		return nil, nil
	}
	return messages, err
}

func (m *defaultChatMessageModel) MarkAsRead(ctx context.Context, messageId int64) error {
	query := fmt.Sprintf("UPDATE %s SET is_read = true WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, messageId)
	return err
}

func (m *defaultChatMessageModel) MarkAsReadBatch(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf("UPDATE %s SET is_read = true WHERE id IN (%s)", m.table, strings.Join(placeholders, ","))
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func (m *defaultChatMessageModel) FindLatestMessageBySession(ctx context.Context, sessionId int64) (*ChatMessage, error) {
	query := fmt.Sprintf(`
    SELECT id, sender_id, receiver_id, content, send_time, is_read
    FROM %s
    WHERE session_id = ?
    ORDER BY send_time DESC
    LIMIT 1
`, m.table)

	var msg ChatMessage
	err := m.conn.QueryRowCtx(ctx, &msg, query, sessionId)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
