package admin

import (
	"context"
	"fmt"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/message/rpc/types/message"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLRejectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLRejectLogic {
	return &ArticleLRejectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ArticleLReject 驳回文章
func (l *ArticleLRejectLogic) ArticleLReject(req *types.ArticleLRejectRequest) (resp *types.ArticleLRejectResponse, err error) {
	_, err = l.svcCtx.ArticleRpc.RejectArticle(l.ctx, &article.ArticleRejectRequest{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("approve article err:%v", err)
		return nil, err
	}
	articleDetail, err := l.svcCtx.ArticleRpc.ArticleDetail(l.ctx, &article.ArticleDetailRequest{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get article detail err:%v", err)
		return nil, err
	}

	_, err = l.svcCtx.MessageRpc.SendMessage(l.ctx, &message.SendMessageRequest{
		BizId:      "article",
		TargetId:   req.ArticleId,
		Title:      "文章审核通知",
		Content:    fmt.Sprintf("文章审核不通过:%s", req.Reason),
		ReceiverId: articleDetail.Article.AuthorId,
		Type:       0,
	})
	if err != nil {
		logx.Errorf("send message err:%v", err)
	}
	return &types.ArticleLRejectResponse{}, nil
}
