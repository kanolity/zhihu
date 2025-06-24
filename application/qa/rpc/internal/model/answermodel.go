package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AnswerModel = (*customAnswerModel)(nil)

type (
	// AnswerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAnswerModel.
	AnswerModel interface {
		answerModel
		withSession(session sqlx.Session) AnswerModel
		ListByQuestion(ctx context.Context, questionId, cursor, limit int64) ([]*Answer, error)
	}

	customAnswerModel struct {
		*defaultAnswerModel
	}
)

// NewAnswerModel returns a model for the database table.
func NewAnswerModel(conn sqlx.SqlConn) AnswerModel {
	return &customAnswerModel{
		defaultAnswerModel: newAnswerModel(conn),
	}
}

func (m *customAnswerModel) withSession(session sqlx.Session) AnswerModel {
	return NewAnswerModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultAnswerModel) ListByQuestion(ctx context.Context, questionId, cursor, limit int64) ([]*Answer, error) {
	var args []interface{}
	var query string

	if cursor > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE question_id = ? AND id < ? ORDER BY id DESC LIMIT ?", answerRows, m.table)
		args = []interface{}{questionId, cursor, limit}
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE question_id = ? ORDER BY id DESC LIMIT ?", answerRows, m.table)
		args = []interface{}{questionId, limit}
	}

	var answers []*Answer
	err := m.conn.QueryRowsCtx(ctx, &answers, query, args...)
	return answers, err
}
