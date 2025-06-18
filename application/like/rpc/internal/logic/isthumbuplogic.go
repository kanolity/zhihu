package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &like.IsThumbupResponse{}, nil
}
