package logic

import (
	"context"
	"go_code/zhihu/application/like/api/internal/code"
	"go_code/zhihu/application/like/rpc/likeclient"

	"go_code/zhihu/application/like/api/internal/svc"
	"go_code/zhihu/application/like/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsThumbupLogic) IsThumbup(req *types.IsThumbupReq) (resp *types.IsThumbupResp, err error) {
	response, err := l.svcCtx.LikeRpc.IsThumbup(l.ctx, &likeclient.IsThumbupRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, code.GetIsThumbupFailed
	}

	// 转换 RPC 响应 到 API 响应
	data := make(map[int64]types.UserThumbup)
	for k, v := range response.UserThumbups {
		data[k] = types.UserThumbup{
			UserId:      v.UserId,
			ThumbupTime: v.ThumbupTime,
			LikeType:    v.LikeType,
		}
	}

	return &types.IsThumbupResp{UserThumbups: data}, nil
}
