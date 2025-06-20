package logic

import (
	"context"
	"go_code/zhihu/application/article/rpc/internal/code"
	"go_code/zhihu/application/article/rpc/internal/model"
	"time"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *article.PublishRequest) (*article.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		Title:       in.Title,
		Content:     in.Content,
		AuthorId:    uint64(in.UserId),
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		logx.Errorf("publish insert req:%v err: %v", in, err)
		return nil, err
	}
	articleId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("publish last insert id err:%v", err)
		return nil, err
	}
	return &article.PublishResponse{ArticleId: articleId}, nil
}
