package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"go_code/zhihu/application/article/mq/internal/svc"
	"go_code/zhihu/application/article/mq/internal/types"
	"strconv"
)

type ArticleLikeNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLikeNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeNumLogic {
	return &ArticleLikeNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLikeNumLogic) Consume(ctx context.Context, key string, val string) error {
	fmt.Println("原始数据：", val)

	var msg *types.CanalLikeMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}
	return l.handleLikeUpdate(ctx, msg)
}

func (l *ArticleLikeNumLogic) handleLikeUpdate(ctx context.Context, msg *types.CanalLikeMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}

	for _, d := range msg.Data {
		if d.BizID != types.ArticleBizID {
			continue
		}

		articleID, err := strconv.ParseInt(d.TargetId, 10, 64)
		if err != nil {
			logx.Errorf("解析文章 ID 错误: %s, err: %v", d.TargetId, err)
			continue
		}

		likeNum, err := strconv.ParseInt(d.LikeNum, 10, 64)
		if err != nil {
			logx.Errorf("解析点赞数错误: %s, err: %v", d.LikeNum, err)
			continue
		}

		if changed, err := l.shouldUpdateLikeNum(ctx, articleID, likeNum); err != nil {
			logx.Errorf("检查是否需要更新失败, id=%d: %v", articleID, err)
			continue
		} else if !changed {
			continue
		}

		// 真正执行更新
		err = l.svcCtx.ArticleModel.UpdateLikeNum(ctx, articleID, likeNum)
		if err != nil {
			logx.Errorf("更新点赞数失败 id=%d like=%d: %v", articleID, likeNum, err)
			continue
		}
	}
	return nil
}

func (l *ArticleLikeNumLogic) shouldUpdateLikeNum(ctx context.Context, id int64, newLikeNum int64) (bool, error) {
	article, err := l.svcCtx.ArticleModel.FindOne(ctx, id)
	if err != nil {
		return false, err
	}
	if article.LikeNum == newLikeNum {
		return false, nil // 没变化就跳过
	}
	return true, nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewArticleLikeNumLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumerConf, NewArticleLogic(ctx, svcCtx)),
	}
}
