package logic

import (
	"context"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"go_code/zhihu/application/reply/api/internal/svc"
	"go_code/zhihu/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRepliesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepliesLogic {
	return &GetRepliesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRepliesLogic) GetReplies(req *types.GetRepliesReq) (resp *types.GetRepliesResp, err error) {
	response, err := l.svcCtx.ReplyRpc.GetReplies(l.ctx, &reply.GetRepliesRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
	})
	if err != nil {
		return nil, err
	}

	replies := make([]types.Reply, 0, len(response.Replies))
	for _, r := range response.Replies {
		replies = append(replies, types.Reply{
			Id:            r.Id,
			BizId:         r.BizId,
			TargetId:      r.TargetId,
			ReplyUserId:   r.ReplyUserId,
			BeReplyUserId: r.BeReplyUserId,
			ParentId:      r.ParentId,
			Content:       r.Content,
			LikeNum:       r.LikeNum,
			IsDeleted:     r.IsDeleted,
			CreateTime:    r.CreateTime,
		})
	}

	return &types.GetRepliesResp{
		Replies: replies,
		HasMore: response.HasMore,
	}, nil
}
