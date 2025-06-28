package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/article/rpc/types/article"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeletedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDeletedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeletedLogic {
	return &ArticleDeletedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDeletedLogic) ArticleDeleted(req *types.ArticleDeletedRequest) (resp *types.ArticleDeletedResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("get userid from context err:%v", err)
		return nil, err
	}
	_, err = l.svcCtx.ArticleRpc.ArticleDelete(l.ctx, &article.ArticleDeleteRequest{
		ArticleId: req.ArticleId,
		UserId:    userId,
	})
	if err != nil {
		logx.Errorf("delete article err:%v", err)
		return nil, err
	}
	return &types.ArticleDeletedResponse{}, nil
}
