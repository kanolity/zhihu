package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"go_code/zhihu/application/like/api/internal/code"
	"go_code/zhihu/application/like/api/internal/svc"
	"go_code/zhihu/application/like/api/internal/types"
	"go_code/zhihu/application/like/rpc/likeclient"
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

	//if req.LikeType == 0 {
	//	l.svcCtx.MessageRpc.SendMessage(l.ctx, &messageservice.SendMessageRequest{
	//		Type:       2,
	//		BizId:      "like",
	//		TargetId:   req.ObjId,
	//		Title:      "收到点赞",
	//		Content:    "收到点赞",
	//		ReceiverId: req.UserId,
	//	})
	//}

	return &types.ThumbupResp{
		BizId:      response.BizId,
		ObjId:      response.ObjId,
		LikeNum:    response.LikeNum,
		DislikeNum: response.DislikeNum,
	}, nil
}
