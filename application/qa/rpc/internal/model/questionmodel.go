package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ QuestionModel = (*customQuestionModel)(nil)

type (
	// QuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionModel.
	QuestionModel interface {
		questionModel
		withSession(session sqlx.Session) QuestionModel
		ListQuestionsByCursor(ctx context.Context, cursorTime int64, lastId int64, limit int64) ([]*Question, error)
	}

	customQuestionModel struct {
		*defaultQuestionModel
	}
)

// NewQuestionModel returns a model for the database table.
func NewQuestionModel(conn sqlx.SqlConn) QuestionModel {
	return &customQuestionModel{
		defaultQuestionModel: newQuestionModel(conn),
	}
}

func (m *customQuestionModel) withSession(session sqlx.Session) QuestionModel {
	return NewQuestionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultQuestionModel) ListQuestionsByCursor(ctx context.Context, cursorTime int64, lastId int64, limit int64) ([]*Question, error) {
	var (
		questions []*Question
		query     string
		args      []any
	)

	query = `
            SELECT id, user_id, title, description, is_resolved, create_time,update_time
            FROM question
            WHERE (UNIX_TIMESTAMP(create_time)  < ?)
               OR (UNIX_TIMESTAMP(create_time)  = ? AND id < ?)
            ORDER BY create_time DESC, id DESC
            LIMIT ?
        `
	args = []any{cursorTime, cursorTime, lastId, limit}

	err := m.conn.QueryRowsCtx(ctx, &questions, query, args...)
	if err != nil {
		return nil, err
	}
	return questions, nil
}
