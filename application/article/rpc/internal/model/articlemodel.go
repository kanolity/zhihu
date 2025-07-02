package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		FindPendingArticlesByStatusWithCursor(ctx context.Context, status int, cursorTime string, lastId int64, limit int64) ([]*Article, error)
		IncreaseReplyNum(ctx context.Context, articleId int64) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}
func (m *customArticleModel) ArticlesByUserId(ctx context.Context, userId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Article, error) {
	var (
		err      error
		sql1     string
		anyField any
		articles []*Article
	)
	if sortField == "like_num" {
		anyField = likeNum
		sql1 = fmt.Sprintf("select "+articleRows+" from "+m.table+" where author_id=? and status=? and like_num < ? order by %s desc limit ?", sortField)
	} else {
		anyField = pubTime
		sql1 = fmt.Sprintf("select "+articleRows+" from "+m.table+" where author_id=? and status=? and publish_time < ? order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &articles, sql1, userId, status, anyField, limit)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (m *customArticleModel) UpdateArticleStatus(ctx context.Context, id int64, status int) error {
	beyondArticleArticleIdKey := fmt.Sprintf("%s%v", cacheBeyondArticleArticleIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set status = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, beyondArticleArticleIdKey)

	return err
}
func (m *customArticleModel) FindPendingArticlesByStatusWithCursor(
	ctx context.Context,
	status int,
	cursorTime string,
	lastId int64,
	limit int64,
) ([]*Article, error) {
	var query string
	var articles []*Article
	var err error
	if lastId == 0 {
		query = fmt.Sprintf(`
        SELECT %s FROM %s 
        WHERE status = ? 
        AND UNIX_TIMESTAMP(publish_time) < ?
        ORDER BY publish_time DESC, id DESC 
        LIMIT ?`, articleRows, m.table)
		err = m.QueryRowsNoCacheCtx(ctx, &articles, query, status, cursorTime, limit)
	} else {
		query = fmt.Sprintf(`
        SELECT %s FROM %s 
        WHERE status = ? 
        AND (UNIX_TIMESTAMP(publish_time) < ? OR (UNIX_TIMESTAMP(publish_time) = ? AND id < ?)) 
        ORDER BY publish_time DESC, id DESC 
        LIMIT ?`, articleRows, m.table)
		err = m.QueryRowsNoCacheCtx(ctx, &articles, query, status, cursorTime, cursorTime, lastId, limit)
	}

	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *defaultArticleModel) IncreaseReplyNum(ctx context.Context, articleId int64) error {
	beyondArticleArticleIdKey := fmt.Sprintf("%s%v", cacheBeyondArticleArticleIdPrefix, articleId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set comment_num = comment_num + 1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, articleId)
	}, beyondArticleArticleIdKey)
	return err
}
