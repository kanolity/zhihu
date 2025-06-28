package logic

import (
	"context"
	"fmt"
	"go_code/zhihu/application/article/rpc/types"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectArticleLogic {
	return &RejectArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RejectArticle 后台：审核驳回
func (l *RejectArticleLogic) RejectArticle(in *article.ArticleRejectRequest) (*article.ArticleRejectResponse, error) {
	article1, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		logx.Errorf("get article err:%v", err)
		return nil, err
	}
	if article1.Status == types.ArticleStatusVisible || article1.Status == types.ArticleStatusUserDelete {
		return nil, fmt.Errorf("文章当前状态无法审核")
	}
	article1.Status = types.ArticleStatusNotPass
	err = l.svcCtx.ArticleModel.Update(l.ctx, article1)
	if err != nil {
		logx.Errorf("reject article err:%v", err)
		return nil, err
	}
	return &article.ArticleRejectResponse{}, nil
}
