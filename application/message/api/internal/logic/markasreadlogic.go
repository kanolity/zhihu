package logic

import (
	"context"
	"go_code/zhihu/application/message/rpc/types/message"

	"go_code/zhihu/application/message/api/internal/svc"
	"go_code/zhihu/application/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAsReadLogic {
	return &MarkAsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAsReadLogic) MarkAsRead(req *types.MarkAsReadReq) (resp *types.MarkAsReadResp, err error) {
	_, err = l.svcCtx.MessageRpc.MarkAsRead(l.ctx, &message.MarkAsReadRequest{
		Id: req.Id,
	})
	return &types.MarkAsReadResp{}, err
}
