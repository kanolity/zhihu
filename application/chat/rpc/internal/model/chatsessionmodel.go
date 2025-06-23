package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatSessionModel = (*customChatSessionModel)(nil)

type (
	// ChatSessionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatSessionModel.
	ChatSessionModel interface {
		chatSessionModel
		withSession(session sqlx.Session) ChatSessionModel
		FindByUserPair(ctx context.Context, user1Id, user2Id int64) (*ChatSession, error)
	}

	customChatSessionModel struct {
		*defaultChatSessionModel
	}
)

// NewChatSessionModel returns a model for the database table.
func NewChatSessionModel(conn sqlx.SqlConn) ChatSessionModel {
	return &customChatSessionModel{
		defaultChatSessionModel: newChatSessionModel(conn),
	}
}

func (m *customChatSessionModel) withSession(session sqlx.Session) ChatSessionModel {
	return NewChatSessionModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultChatSessionModel) FindByUserPair(ctx context.Context, user1Id, user2Id int64) (*ChatSession, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?) LIMIT 1", chatSessionRows, m.table)
	var resp ChatSession
	err := m.conn.QueryRowCtx(ctx, &resp, query, user1Id, user2Id, user2Id, user1Id)
	if err == sqlc.ErrNotFound {
		return nil, nil
	}
	return &resp, err
}
