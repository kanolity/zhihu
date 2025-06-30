package model

import (
	"context"
	"fmt"
	sqlx2 "github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TagModel = (*customTagModel)(nil)

type (
	// TagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagModel.
	TagModel interface {
		tagModel
		withSession(session sqlx.Session) TagModel
		ListTags(ctx context.Context, cursor, limit int64) ([]*Tag, error)
		BatchGetTags(ctx context.Context, ids []int64) ([]*Tag, error)
		FindNamesByIds(ctx context.Context, ids []int64) ([]*Tag, error)
	}

	customTagModel struct {
		*defaultTagModel
	}
)

// NewTagModel returns a model for the database table.
func NewTagModel(conn sqlx.SqlConn) TagModel {
	return &customTagModel{
		defaultTagModel: newTagModel(conn),
	}
}

func (m *customTagModel) withSession(session sqlx.Session) TagModel {
	return NewTagModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultTagModel) ListTags(ctx context.Context, cursor, limit int64) ([]*Tag, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id > ? ORDER BY id DESC LIMIT ?`, tagRows, m.table)
	var tags []*Tag
	err := m.conn.QueryRowsCtx(ctx, &tags, query, cursor, limit)
	logx.Infof("ListTags: cursor=%d, limit=%d, got=%d tags", cursor, limit, len(tags))

	return tags, err
}

func (m *defaultTagModel) BatchGetTags(ctx context.Context, ids []int64) ([]*Tag, error) {
	if len(ids) == 0 {
		return []*Tag{}, nil
	}
	inQuery, inArgs, _ := sqlx2.In("SELECT "+tagRows+" FROM "+m.table+" WHERE id IN (?)", ids)

	var tags []*Tag
	err := m.conn.QueryRowsCtx(ctx, &tags, inQuery, inArgs...)
	return tags, err
}

func (m *defaultTagModel) FindNamesByIds(ctx context.Context, ids []int64) ([]*Tag, error) {
	if len(ids) == 0 {
		return []*Tag{}, nil
	}

	query, args, err := sqlx2.In(fmt.Sprintf("SELECT %s FROM %s WHERE id IN (?)", tagRows, m.table), ids)
	if err != nil {
		return nil, err
	}

	var tags []*Tag
	err = m.conn.QueryRowsCtx(ctx, &tags, query, args...)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
