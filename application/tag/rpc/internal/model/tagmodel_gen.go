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
	tagFieldNames          = builder.RawFieldNames(&Tag{})
	tagRows                = strings.Join(tagFieldNames, ",")
	tagRowsExpectAutoSet   = strings.Join(stringx.Remove(tagFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tagRowsWithPlaceHolder = strings.Join(stringx.Remove(tagFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tagModel interface {
		Insert(ctx context.Context, data *Tag) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Tag, error)
		Update(ctx context.Context, data *Tag) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTagModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Tag struct {
		Id         uint64    `db:"id"`          // 主键ID
		TagName    string    `db:"tag_name"`    // 标签名
		TagDesc    string    `db:"tag_desc"`    // 标签描述
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 最后修改时间
	}
)

func newTagModel(conn sqlx.SqlConn) *defaultTagModel {
	return &defaultTagModel{
		conn:  conn,
		table: "`tag`",
	}
}

func (m *defaultTagModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTagModel) FindOne(ctx context.Context, id uint64) (*Tag, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tagRows, m.table)
	var resp Tag
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

func (m *defaultTagModel) Insert(ctx context.Context, data *Tag) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tagRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.TagName, data.TagDesc)
	return ret, err
}

func (m *defaultTagModel) Update(ctx context.Context, data *Tag) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tagRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.TagName, data.TagDesc, data.Id)
	return err
}

func (m *defaultTagModel) tableName() string {
	return m.table
}
