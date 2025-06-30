package logic

import (
	"context"
	"go_code/zhihu/application/article/rpc/internal/code"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleReplyIncreaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleReplyIncreaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleReplyIncreaseLogic {
	return &ArticleReplyIncreaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleReplyIncreaseLogic) ArticleReplyIncrease(in *article.ArticleReplyIncreaseRequest) (*article.ArticleReplyIncreaseResponse, error) {
	if in.ArticleId <= 0 {
		return nil, code.ArticleIdInvalid
	}

	err := l.svcCtx.ArticleModel.IncreaseReplyNum(l.ctx, in.ArticleId)
	if err != nil {
		return nil, err
	}
	return &article.ArticleReplyIncreaseResponse{}, nil
}
