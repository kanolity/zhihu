package logic

import (
	"context"

	"go_code/zhihu/application/like/rpc/internal/svc"
	"go_code/zhihu/application/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *like.ThumbupRequest) (*like.ThumbupResponse, error) {
	// todo: add your logic here and delete this line

	return &like.ThumbupResponse{}, nil
}
