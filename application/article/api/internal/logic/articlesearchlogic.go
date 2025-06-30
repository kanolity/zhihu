package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleSearchLogic {
	return &ArticleSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleSearchLogic) ArticleSearch(req *types.ArticleSearchRequest) (resp *types.ArticleSearchResponse, err error) {
	resp = &types.ArticleSearchResponse{
		Articles: make([]types.ESArticleInfo, 0),
	}

	// 设置默认分页大小
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建 ES 查询体
	query := map[string]interface{}{}

	// 构建 bool 查询
	boolQuery := map[string]interface{}{}

	// 关键词搜索 multi_match 查询，支持标题、内容、描述等字段
	if req.Query != "" {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     req.Query,
					"fields":    []string{"title", "content", "description"},
					"fuzziness": 2,
				},
			},
		}
	} else {
		boolQuery["must"] = []interface{}{
			map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}

	// 过滤条件
	filterArr := make([]interface{}, 0)
	//filterArr = append(filterArr, map[string]interface{}{
	//	"term": map[string]interface{}{
	//		"status": 2,
	//	},
	//})

	// 作者过滤
	if req.AuthorId > 0 {
		filterArr = append(filterArr, map[string]interface{}{
			"term": map[string]interface{}{
				"author_id": req.AuthorId,
			},
		})
	}
	// 标签过滤（term 或 terms 查询）
	if len(req.TagNames) > 0 {
		filterArr = append(filterArr, map[string]interface{}{
			"terms": map[string]interface{}{
				"tag_names": req.TagNames,
			},
		})
	}

	if len(filterArr) > 0 {
		boolQuery["filter"] = filterArr
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	// 排序字段及默认排序
	sortField := "publish_time"
	if req.SortType == 1 {
		sortField = "like_num"
	}

	query["sort"] = []interface{}{
		map[string]interface{}{sortField: map[string]string{"order": "desc"}},
		map[string]interface{}{"article_id": map[string]string{"order": "desc"}}, // 防止排序字段相同时，ID作为二次排序
	}

	// 分页，基于 search_after 游标分页
	if req.Cursor > 0 && req.ArticleId > 0 {
		query["search_after"] = []interface{}{req.Cursor, req.ArticleId}
	}

	query["size"] = req.PageSize

	// 编码成 JSON 请求体
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logx.Errorf("json encode err:%v", err)
		return nil, err
	}

	// 调用 ES 搜索接口
	res, err := l.svcCtx.Es.Client.Search(
		l.svcCtx.Es.Client.Search.WithContext(l.ctx),
		l.svcCtx.Es.Client.Search.WithIndex("article-index"),
		l.svcCtx.Es.Client.Search.WithBody(strings.NewReader(buf.String())),
	)

	if err != nil {
		logx.Errorf("es client err:%v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	logx.Infof("ES Raw Response: %s", string(body))

	if res.IsError() {
		logx.Errorf("elasticsearch search err: %s", string(body))
		return nil, fmt.Errorf("elasticsearch search error: %s", string(body))
	}

	// 解析响应
	var r struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source types.ESArticleInfo `json:"_source"`
				Sort   []interface{}       `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &r); err != nil {
		logx.Errorf("decode response body err: %v", err)
		return nil, err
	}

	resp.Total = r.Hits.Total.Value

	for _, hit := range r.Hits.Hits {
		resp.Articles = append(resp.Articles, hit.Source)

		// 更新游标为最后一条的排序字段和文章ID
		if len(hit.Sort) >= 2 {
			if val, ok := hit.Sort[0].(float64); ok {
				resp.Cursor = int64(val)
			}
			if val, ok := hit.Sort[1].(float64); ok {
				resp.ArticleId = int64(val)
			}
		}
	}

	resp.IsEnd = len(r.Hits.Hits) < int(req.PageSize)

	return resp, nil
}
