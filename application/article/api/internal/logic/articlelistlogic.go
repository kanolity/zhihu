package logic

import (
	"context"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListRequest) (resp *types.ArticleListResponse, err error) {
	articles, err := l.svcCtx.ArticleRpc.Articles(l.ctx, &article.ArticlesRequest{
		UserId:    req.AuthorId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get articles req: %v err: %v", req, err)
		return nil, err
	}
	if articles == nil || len(articles.Articles) == 0 {
		return &types.ArticleListResponse{}, nil
	}
	articleInfos := make([]types.ArticleInfo, 0, len(articles.Articles))
	for _, a := range articles.Articles {
		tags, err := l.svcCtx.TagRpc.GetTags(l.ctx, &tag.GetTagsRequest{TagIds: a.TagIds})
		if err != nil {
			logx.Errorf("get tags req: %v err: %v", req, err)
			continue
		}
		articleInfos = append(articleInfos, types.ArticleInfo{
			ArticleId: a.Id,
			Title:     a.Title,
			TagNames:  tags.TagNames,
			LikeNum:   a.LikeCount,
		})
	}

	return &types.ArticleListResponse{
		Articles: articleInfos,
	}, nil
}
