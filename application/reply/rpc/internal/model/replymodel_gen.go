// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.8.3

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	replyFieldNames          = builder.RawFieldNames(&Reply{})
	replyRows                = strings.Join(replyFieldNames, ",")
	replyRowsExpectAutoSet   = strings.Join(stringx.Remove(replyFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	replyRowsWithPlaceHolder = strings.Join(stringx.Remove(replyFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	replyModel interface {
		Insert(ctx context.Context, data *Reply) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Reply, error)
		Update(ctx context.Context, data *Reply) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultReplyModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Reply struct {
		Id            uint64    `db:"id"`               // 主键ID
		BizId         string    `db:"biz_id"`           // 业务ID
		TargetId      uint64    `db:"target_id"`        // 评论目标id
		ReplyUserId   uint64    `db:"reply_user_id"`    // 评论用户ID
		BeReplyUserId uint64    `db:"be_reply_user_id"` // 被回复用户ID
		ParentId      uint64    `db:"parent_id"`        // 父评论ID
		Content       string    `db:"content"`          // 内容
		Status        int64     `db:"status"`           // 状态 0:正常 1:删除
		LikeNum       int64     `db:"like_num"`         // 点赞数
		CreateTime    time.Time `db:"create_time"`      // 创建时间
		UpdateTime    time.Time `db:"update_time"`      // 最后修改时间
	}
)

func newReplyModel(conn sqlx.SqlConn) *defaultReplyModel {
	return &defaultReplyModel{
		conn:  conn,
		table: "`reply`",
	}
}

func (m *defaultReplyModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultReplyModel) FindOne(ctx context.Context, id uint64) (*Reply, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", replyRows, m.table)
	var resp Reply
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultReplyModel) Insert(ctx context.Context, data *Reply) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, replyRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.BizId, data.TargetId, data.ReplyUserId, data.BeReplyUserId, data.ParentId, data.Content, data.Status, data.LikeNum)
	return ret, err
}

func (m *defaultReplyModel) Update(ctx context.Context, data *Reply) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, replyRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.BizId, data.TargetId, data.ReplyUserId, data.BeReplyUserId, data.ParentId, data.Content, data.Status, data.LikeNum, data.Id)
	return err
}

func (m *defaultReplyModel) tableName() string {
	return m.table
}
