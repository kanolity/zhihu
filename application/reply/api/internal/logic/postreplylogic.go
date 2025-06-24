package logic

import (
	"context"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"go_code/zhihu/application/reply/api/internal/svc"
	"go_code/zhihu/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostReplyLogic {
	return &PostReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostReplyLogic) PostReply(req *types.PostReplyReq) (resp *types.PostReplyResp, err error) {
	response, err := l.svcCtx.ReplyRpc.PostReply(l.ctx, &reply.PostReplyRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   req.ReplyUserId,
		BeReplyUserId: req.BeReplyUserId,
		ParentId:      req.ParentId,
		Content:       req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &types.PostReplyResp{Id: response.Id}, nil
}
