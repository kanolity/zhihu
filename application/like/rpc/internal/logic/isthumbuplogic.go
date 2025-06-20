package logic

import (
	"context"
	"go_code/zhihu/application/like/rpc/internal/model"

	"go_code/zhihu/application/like/rpc/internal/svc"
	"go_code/zhihu/application/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupLogic) IsThumbup(in *like.IsThumbupRequest) (*like.IsThumbupResponse, error) {
	result := make(map[int64]*like.UserThumbup)

	lk, err := l.svcCtx.LikeModel.FindByUnique(l.ctx, in.BizId, in.TargetId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if lk != nil {
		result[int64(lk.UserId)] = &like.UserThumbup{
			UserId:      int64(lk.UserId),
			ThumbupTime: lk.CreateTime.Unix(),
			LikeType:    int32(lk.Type),
		}
	}

	return &like.IsThumbupResponse{UserThumbups: result}, nil
}
