package admin

import (
	"context"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/message/rpc/types/message"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleApproveLogic {
	return &ArticleApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ArticleApprove 审核通过
func (l *ArticleApproveLogic) ArticleApprove(req *types.ArticleApproveRequest) (resp *types.ArticleApproveResponse, err error) {
	_, err = l.svcCtx.ArticleRpc.ApproveArticle(l.ctx, &article.ArticleApproveRequest{
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
		Content:    "文章审核通过",
		ReceiverId: articleDetail.Article.AuthorId,
		Type:       0,
	})
	if err != nil {
		logx.Errorf("send message err:%v", err)
	}
	return &types.ArticleApproveResponse{}, nil
}
