package logic

import (
	"cmp"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"go_code/zhihu/application/article/rpc/internal/code"
	"go_code/zhihu/application/article/rpc/internal/model"
	"go_code/zhihu/application/article/rpc/types"
	"slices"
	"strconv"
	"time"

	"go_code/zhihu/application/article/rpc/internal/svc"
	"go_code/zhihu/application/article/rpc/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixArticles = "biz#articles#%d#%d"
	articlesExpire = 3600 * 24 * 2
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *article.ArticlesRequest) (*article.ArticlesResponse, error) {
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	var (
		sortField       string
		sortLikeNum     int64
		sortPublishTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	}

	var (
		err            error
		isCache, isEnd bool
		lastId, cursor int64
		curPage        []*article.ArticleItem
		articles       []*model.Article
	)
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	if len(articleIds) > 0 {
		isCache = true
		if articleIds[len(articleIds)-1] == -1 {
			isEnd = true
		}
		articles, err = l.articleByIds(l.ctx, articleIds)
		if err != nil {
			return nil, err
		}

		// 通过sortFiled对articles进行排序
		var cmpFunc func(a, b *model.Article) int
		if sortField == "like_num" {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.LikeNum, a.LikeNum)
			}
		} else {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
			}
		}
		slices.SortFunc(articles, cmpFunc)

		for _, article1 := range articles {
			curPage = append(curPage, &article.ArticleItem{
				Id:           article1.Id,
				Title:        article1.Title,
				Content:      article1.Content,
				LikeCount:    article1.LikeNum,
				CommentCount: article1.CommentNum,
				PublishTime:  article1.PublishTime.Unix(),
				TagIds:       article1.TagIds,
			})
		}
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("ArticlesByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.ArticlesByUserId(l.ctx, in.UserId, types.ArticleStatusVisible, sortLikeNum, sortPublishTime, sortField, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, sortField, err)
			return nil, err
		}
		if v == nil {
			return &article.ArticlesResponse{}, nil
		}
		articles = v.([]*model.Article)
		var firstPageArticles []*model.Article
		if len(articles) > int(in.PageSize) {
			firstPageArticles = articles[:int(in.PageSize)]
		} else {
			firstPageArticles = articles
			isEnd = true
		}
		for _, article1 := range firstPageArticles {
			curPage = append(curPage, &article.ArticleItem{
				Id:           article1.Id,
				Title:        article1.Title,
				Content:      article1.Content,
				LikeCount:    article1.LikeNum,
				CommentCount: article1.CommentNum,
				PublishTime:  article1.PublishTime.Unix(),
				TagIds:       article1.TagIds,
			})
		}
	}

	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		if in.SortType == types.SortPublishTime {
			cursor = pageLast.PublishTime
		} else {
			cursor = pageLast.LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, article1 := range curPage {
			if in.SortType == types.SortPublishTime {
				if article1.PublishTime == in.Cursor && article1.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			} else {
				if article1.LikeCount == in.Cursor && article1.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	ret := &article.ArticlesResponse{
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
		Articles:  curPage,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(articles) < types.DefaultLimit && len(articles) > 0 {
				articles = append(articles, &model.Article{Id: -1})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheArticles error: %v", err)
			}
		})
	}

	return ret, nil
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixArticles, uid, sortType)
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, uid, cursor, ps int64, sortType int32) ([]int64, error) {
	key := articlesKey(uid, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimit key: %s error: %v", key, err)
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, articles []*model.Article, userId int64, sortType int32) error {
	if len(articles) == 0 {
		return nil
	}
	key := articlesKey(userId, sortType)
	for _, article1 := range articles {
		var score int64
		if sortType == types.SortLikeCount {
			score = article1.LikeNum
		} else if sortType == types.SortPublishTime && article1.Id != -1 {
			score = article1.PublishTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(article1.Id)))
		if err != nil {
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
}

func (l *ArticlesLogic) articleByIds(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](func(source chan<- int64) {
		for _, aid := range articleIds {
			if aid == -1 {
				continue
			}
			source <- aid
		}
	}, func(id int64, writer mr.Writer[*model.Article], cancel func(error)) {
		p, err := l.svcCtx.ArticleModel.FindOne(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
		var articles []*model.Article
		for article1 := range pipe {
			articles = append(articles, article1)
		}
		writer.Write(articles)
	})
	if err != nil {
		return nil, err
	}

	return articles, nil
}
