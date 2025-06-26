package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindUsersByIds(ctx context.Context, ids []int64) ([]*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}
func (m *defaultUserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.QueryRowNoCacheCtx(ctx, &user,
		fmt.Sprintf("select %s from %s where `username` = ? limit 1", userRows, m.table), username)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *defaultUserModel) FindByMobile(ctx context.Context, mobile string) (*User, error) {
	var user User
	err := m.QueryRowNoCacheCtx(ctx, &user,
		fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", userRows, m.table), mobile)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *defaultUserModel) FindUsersByIds(ctx context.Context, ids []int64) ([]*User, error) {
	if len(ids) == 0 {
		return []*User{}, nil
	}

	users := make([]*User, 0, len(ids))
	for _, id := range ids {
		var u User
		cacheKey := fmt.Sprintf("%s%d", cacheBeyondUserUserIdPrefix, id)

		err := m.CachedConn.QueryRowCtx(ctx, &u, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
			query := fmt.Sprintf("SELECT id, username, avatar FROM %s WHERE id = ?", m.table)
			return conn.QueryRowCtx(ctx, v, query, id)
		})

		if err != nil {
			if err != sqlc.ErrNotFound {
				logx.Errorf("FindUsersByIds: id=%d error=%v", id, err)
			}
			continue // 未命中或错误则跳过
		}
		users = append(users, &u)
	}
	return users, nil
}
