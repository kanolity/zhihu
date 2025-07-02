package admin

import (
	"context"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/tag/rpc/types/tag"
	"go_code/zhihu/application/user/rpc/types/user"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListLogic1 struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListLogic1 {
	return &AdminListLogic1{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AdminList 待审核列表
func (l *AdminListLogic1) AdminList(req *types.AdminListRequest) (resp *types.AdminListResponse, err error) {
	response, err := l.svcCtx.ArticleRpc.GetPendingArticles(l.ctx, &article.AdminListRequest{
		ArticleId: req.ArticleId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		Status:    req.Status,
	})
	if err != nil {
		logx.Errorf("Get PendingArticles err:%v", err)
		return nil, err
	}

	var items []types.PendingArticleItem
	for _, a := range response.Articles {
		tags, err := l.svcCtx.TagRpc.GetTags(l.ctx, &tag.GetTagsRequest{
			TagIds: a.TagIds,
		})
		if err != nil {
			logx.Errorf("Get Tags err:%v", err)
		}
		authorName, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
			UserId: a.AuthorId,
		})
		if err != nil {
			logx.Errorf("Get Author err:%v", err)
		}

		items = append(items, types.PendingArticleItem{
			ArticleId:  a.ArticleId,
			Title:      a.Title,
			AuthorId:   a.AuthorId,
			AuthorName: authorName.Username,
			TagNames:   tags.TagNames,
			Status:     a.Status,
		})
	}
	return &types.AdminListResponse{
		Articles:  items,
		Cursor:    response.Cursor,
		ArticleId: response.ArticleId,
		IsEnd:     response.IsEnd,
	}, nil
}
