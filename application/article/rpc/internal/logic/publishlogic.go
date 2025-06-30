package logic

import (
	"context"
	"fmt"
	"go_code/zhihu/application/article/rpc/internal/code"
	"go_code/zhihu/application/article/rpc/internal/model"
	"go_code/zhihu/application/article/rpc/types"
	"strconv"
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
	fmt.Println("PublishLogic Publish")
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		Title:       in.Title,
		Content:     in.Content,
		AuthorId:    uint64(in.UserId),
		Status:      types.ArticleStatusPending,
		TagIds:      in.TagIds,
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
	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
	)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, time.Now().Unix(), articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
	}
	b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
	}
	return &article.PublishResponse{ArticleId: articleId}, nil
}
