package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/article/rpc/internal/code"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *article.ArticleDetailRequest) (*article.ArticleDetailResponse, error) {
	article1, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &article.ArticleDetailResponse{}, nil
		}
		return nil, code.FindArticleFailed
	}
	return &article.ArticleDetailResponse{
		Article: &article.ArticleItem{
			Id:          article1.Id,
			Title:       article1.Title,
			Content:     article1.Content,
			AuthorId:    int64(article1.AuthorId),
			LikeCount:   article1.LikeNum,
			PublishTime: article1.PublishTime.Unix(),
			TagIds:      article1.TagIds,
		},
	}, nil
}
