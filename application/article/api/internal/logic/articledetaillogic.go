package logic

import (
	"context"
	"fmt"
	"go_code/zhihu/application/article/api/internal/code"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/tag/rpc/types/tag"
	"go_code/zhihu/application/user/rpc/types/user"
	"strconv"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDetailLogic) ArticleDetail(req *types.ArticleDetailRequest) (resp *types.ArticleDetailResponse, err error) {
	articleInfo, err := l.svcCtx.ArticleRpc.ArticleDetail(l.ctx, &article.ArticleDetailRequest{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get article detail id: %d err: %v", req.ArticleId, err)
		return nil, code.GetArticleDetailFailed
	}
	if articleInfo == nil || articleInfo.Article == nil {
		return nil, nil
	}
	userInfo, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
		UserId: articleInfo.Article.AuthorId,
	})
	if err != nil {
		logx.Errorf("get userInfo id: %d err: %v", articleInfo.Article.AuthorId, err)
		return nil, code.GetUserInfoFailed
	}

	fmt.Printf("tagids:%+v\n", articleInfo.Article)
	tags, err := l.svcCtx.TagRpc.GetTags(l.ctx, &tag.GetTagsRequest{TagIds: articleInfo.Article.TagIds})
	if err != nil {
		logx.Errorf("get tags req: %v err: %v", req, err)
	}
	return &types.ArticleDetailResponse{
		Title:      articleInfo.Article.Title,
		Content:    articleInfo.Article.Content,
		AuthorId:   strconv.FormatInt(articleInfo.Article.AuthorId, 10),
		AuthorName: userInfo.Username,
		TagNames:   tags.TagNames,
	}, nil
}
