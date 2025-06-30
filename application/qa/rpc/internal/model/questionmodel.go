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
	query := `
        SELECT id, user_id, title, description, is_resolved, create_time
        FROM question
        WHERE (UNIX_TIMESTAMP(create_time) * 1000 < ?) 
   OR (UNIX_TIMESTAMP(create_time) * 1000 = ? AND id < ?)
        ORDER BY create_time DESC, id DESC
        LIMIT ?
    `
	var questions []*Question
	err := m.conn.QueryRowsCtx(ctx, &questions, query, cursorTime, cursorTime, lastId, limit)
	return questions, err
}
