package logic

import (
	"context"
	"encoding/json"
	"go_code/zhihu/application/article/api/internal/code"
	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

const minContentLen = 100

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	if len(req.Title) == 0 {
		return nil, code.TitleCantEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleContentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.CoverCantEmpty
	}
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return nil, xcode.NoLogin
	}
	pRes, err := l.svcCtx.ArticleRpc.Publish(l.ctx, &article.PublishRequest{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		Cover:   req.Cover,
	})
	if err != nil {
		logx.Errorf("l.svcCtx.ArticleRpc.Publish req:%v userId:%v error: %v", req, userId, err)
		return nil, err
	}

	return &types.PublishResponse{ArticleId: pRes.ArticleId}, nil
}
