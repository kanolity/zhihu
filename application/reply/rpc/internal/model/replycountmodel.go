package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReplyCountModel = (*customReplyCountModel)(nil)

type (
	// ReplyCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReplyCountModel.
	ReplyCountModel interface {
		replyCountModel
		withSession(session sqlx.Session) ReplyCountModel
		IncreaseCount(ctx context.Context, bizId string, targetId int64, isRoot bool) error
		FindByArticleId(ctx context.Context, articleId int64) (int64, error)
	}

	customReplyCountModel struct {
		*defaultReplyCountModel
	}
)

// NewReplyCountModel returns a model for the database table.
func NewReplyCountModel(conn sqlx.SqlConn) ReplyCountModel {
	return &customReplyCountModel{
		defaultReplyCountModel: newReplyCountModel(conn),
	}
}

func (m *customReplyCountModel) withSession(session sqlx.Session) ReplyCountModel {
	return NewReplyCountModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultReplyCountModel) IncreaseCount(ctx context.Context, bizId string, targetId int64, isRoot bool) error {
	var query string
	var args []interface{}

	if isRoot {
		query = fmt.Sprintf(`
            INSERT INTO %s (biz_id, target_id, reply_num, reply_root_num)
            VALUES (?, ?, 1, 1)
            ON DUPLICATE KEY UPDATE
                reply_num = reply_num + 1,
                reply_root_num = reply_root_num + 1
        `, m.table)
		args = []interface{}{bizId, targetId}
	} else {
		query = fmt.Sprintf(`
            INSERT INTO %s (biz_id, target_id, reply_num, reply_root_num)
            VALUES (?, ?, 1, 0)
            ON DUPLICATE KEY UPDATE
                reply_num = reply_num + 1
        `, m.table)
		args = []interface{}{bizId, targetId}
	}

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func (m *defaultReplyCountModel) FindByArticleId(ctx context.Context, articleId int64) (int64, error) {
	query := fmt.Sprintf("select %s from %s where `target_id` = ? and `biz_id`=? limit 1", replyCountRows, m.table)
	var resp ReplyCount
	err := m.conn.QueryRowCtx(ctx, &resp, query, articleId, "article")
	switch err {
	case nil:
		return resp.ReplyNum, nil
	case sqlx.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}
