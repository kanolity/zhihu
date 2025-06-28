package logic

import (
	"context"
	"strconv"
	"time"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPendingArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPendingArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPendingArticlesLogic {
	return &GetPendingArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPendingArticles 后台：待审核或驳回文章列表
func (l *GetPendingArticlesLogic) GetPendingArticles(in *article.AdminListRequest) (*article.AdminListResponse, error) {
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	articles, err := l.svcCtx.ArticleModel.FindPendingArticlesByStatusWithCursor(
		l.ctx, int(in.Status), strconv.FormatInt(in.Cursor, 10), in.ArticleId, in.PageSize+1,
	)
	if err != nil {
		logx.Errorf("Find PendingArticles By Status With Cursor err: %v", err)
		return nil, err
	}

	isEnd := true
	if int64(len(articles)) > in.PageSize {
		isEnd = false
		articles = articles[:len(articles)-1]
	}

	var items []*article.PendingArticleItem
	var lastCursor, lastId int64
	for _, a := range articles {
		items = append(items, &article.PendingArticleItem{
			ArticleId:   a.Id,
			Title:       a.Title,
			AuthorId:    int64(a.AuthorId),
			TagIds:      a.TagIds,
			Status:      a.Status,
			PublishTime: a.PublishTime.Unix(),
		})
	}
	if len(items) > 0 {
		last := items[len(items)-1]
		lastCursor = last.PublishTime
		lastId = last.ArticleId
	}

	return &article.AdminListResponse{
		Articles:  items,
		Cursor:    lastCursor,
		ArticleId: lastId,
		IsEnd:     isEnd,
	}, nil
}
