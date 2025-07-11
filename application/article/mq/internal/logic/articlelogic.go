package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/zeromicro/go-zero/core/logx"
	"go_code/zhihu/application/article/mq/internal/svc"
	"go_code/zhihu/application/article/mq/internal/types"
	"go_code/zhihu/application/user/rpc/types/user"
	"strconv"
	"strings"
	"time"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.articleOperate(msg)
}

func (l *ArticleLogic) articleOperate(msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		logx.Infof("msg.Data is empty")
		return nil
	}

	var esData []*types.ArticleEsMsg
	for _, d := range msg.Data {
		status, _ := strconv.Atoi(d.Status)
		likNum, _ := strconv.ParseInt(d.LikeNum, 10, 64)
		fmt.Println("likeNum:", d.LikeNum)
		fmt.Println("likNum:", likNum)
		articleId, _ := strconv.ParseInt(d.ID, 10, 64)
		authorId, _ := strconv.ParseInt(d.AuthorId, 10, 64)

		t, err := time.ParseInLocation("2006-01-02 15:04:05", d.PublishTime, time.Local)
		publishTimeKey := articlesKey(d.AuthorId, 0)
		likeNumKey := articlesKey(d.AuthorId, 1)

		switch status {
		case types.ArticleStatusVisible:
			b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
			if b {
				_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, t.Unix(), d.ID)
				if err != nil {
					l.Logger.Errorf("ZaddCtx key: %s req: %v error: %v", publishTimeKey, d, err)
				}
			}
			b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
			if b {
				_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, likNum, d.ID)
				if err != nil {
					l.Logger.Errorf("ZaddCtx key: %s req: %v error: %v", likeNumKey, d, err)
				}
			}
		case types.ArticleStatusUserDelete:
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, publishTimeKey, d.ID)
			if err != nil {
				l.Logger.Errorf("ZremCtx key: %s req: %v error: %v", publishTimeKey, d, err)
			}
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeNumKey, d.ID)
			if err != nil {
				l.Logger.Errorf("ZremCtx key: %s req: %v error: %v", likeNumKey, d, err)
			}
		}

		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
			UserId: authorId,
		})
		if err != nil {
			l.Logger.Errorf("FindById userId: %d error: %v", authorId, err)
			return err
		}
		publishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", d.PublishTime, time.Local)
		createTime, _ := time.ParseInLocation("2006-01-02 15:04:05", d.CreateTime, time.Local)
		updateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", d.UpdateTime, time.Local)
		cNum, err := strconv.ParseInt(d.CommentNum, 10, 64)
		esData = append(esData, &types.ArticleEsMsg{
			ArticleId:   articleId,
			AuthorId:    authorId,
			AuthorName:  u.Username,
			Title:       d.Title,
			Content:     d.Content,
			Description: d.Description,
			Status:      status,
			LikeNum:     likNum,
			PublishTime: publishTime.Format("2006-01-02 15:04:05"),
			CreateTime:  createTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  updateTime.Format("2006-01-02 15:04:05"),
			CommentNum:  cNum,
		})
	}

	err := l.BatchUpSertToEs(l.ctx, esData)
	logx.Infof("BatchUpSertToEs esData: %v", esData)
	if err != nil {
		l.Logger.Errorf("BatchUpSertToEs data: %v error: %v", esData, err)
	}

	return err
}

func (l *ArticleLogic) BatchUpSertToEs(ctx context.Context, data []*types.ArticleEsMsg) error {
	if len(data) == 0 {
		fmt.Println("msg.Data is empty")
		return nil
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}

	for _, d := range data {
		v, err := json.Marshal(d)
		if err != nil {
			fmt.Println("marsh json err :", err, "docId:", d.ArticleId)
			return err
		}

		payload := fmt.Sprintf(`{"doc":%s,"doc_as_upsert":true}`, string(v))
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "update",
			DocumentID: fmt.Sprintf("%d", d.ArticleId),
			Body:       strings.NewReader(payload),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem) {
				logx.Infof("ES upsert success: docId=%s", item.DocumentID)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem, err error) {
				logx.Errorf("ES upsert fail [ES error]: docId=%s status=%s error=%v", item.DocumentID, item2.Status, item2.Error)
			},
		})
		if err != nil {
			logx.Errorf("bi.Add error for docId=%s: %v", d.ArticleId, err)
			return err
		}
	}

	return bi.Close(ctx)
}

func articlesKey(uid string, sortType int32) string {
	return fmt.Sprintf("biz#articles#%s#%d", uid, sortType)
}
