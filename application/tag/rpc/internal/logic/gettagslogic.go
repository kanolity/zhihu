package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"

	"go_code/zhihu/application/tag/rpc/internal/svc"
	"go_code/zhihu/application/tag/rpc/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagsLogic {
	return &GetTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTagsLogic) GetTags(in *tag.GetTagsRequest) (*tag.GetTagsResponse, error) {
	ids, err := parseTagIDs(in.TagIds)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tag_ids format")
	}

	// 2. 查询数据库：根据 tag_id 批量查询标签名
	tags, err := l.svcCtx.TagModel.FindNamesByIds(l.ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch tag names")
	}

	// 3. 构造返回
	var tagNames []string
	for _, t := range tags {
		tagNames = append(tagNames, t.TagName)
	}

	return &tag.GetTagsResponse{TagNames: tagNames}, nil
}

func parseTagIDs(input string) ([]int64, error) {
	var ids []int64
	for _, str := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(str)
		if trimmed == "" {
			continue
		}
		id, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
