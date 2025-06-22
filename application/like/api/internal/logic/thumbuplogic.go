package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"go_code/zhihu/application/like/api/internal/code"
	"go_code/zhihu/application/like/rpc/likeclient"

	"go_code/zhihu/application/like/api/internal/svc"
	"go_code/zhihu/application/like/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbupLogic) Thumbup(req *types.ThumbupReq) (resp *types.ThumbupResp, err error) {
	response, err := l.svcCtx.LikeRpc.Thumbup(l.ctx, &likeclient.ThumbupRequest{
		BizId:    req.BizId,
		ObjId:    req.ObjId,
		UserId:   req.UserId,
		LikeType: req.LikeType,
	})
	if err != nil && err != sqlc.ErrNotFound {
		return nil, code.ThumbupFailed
	}

	return &types.ThumbupResp{
		BizId:      response.BizId,
		ObjId:      response.ObjId,
		LikeNum:    response.LikeNum,
		DislikeNum: response.DislikeNum,
	}, nil
}
