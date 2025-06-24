package logic

import (
	"context"

	"go_code/zhihu/application/reply/rpc/internal/svc"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRepliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepliesLogic {
	return &GetRepliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRepliesLogic) GetReplies(in *reply.GetRepliesRequest) (*reply.GetRepliesReply, error) {
	list, err := l.svcCtx.ReplyModel.ListByTarget(l.ctx, in.BizId, in.TargetId, in.Cursor, in.Limit+1)
	if err != nil {
		return nil, err
	}

	resp := &reply.GetRepliesReply{
		Replies: make([]*reply.Reply, 0, in.Limit),
		HasMore: false,
	}

	for i, r := range list {
		if int64(i) < in.Limit {
			resp.Replies = append(resp.Replies, &reply.Reply{
				Id:            int64(r.Id),
				BizId:         r.BizId,
				TargetId:      int64(r.TargetId),
				ReplyUserId:   int64(r.ReplyUserId),
				BeReplyUserId: int64(r.BeReplyUserId),
				ParentId:      int64(r.ParentId),
				Content:       r.Content,
				LikeNum:       r.LikeNum,
				IsDeleted:     r.Status != 0,
				CreateTime:    r.CreateTime.Format("2006-01-02 15:04:05"),
			})
		} else {
			resp.HasMore = true
		}
	}

	return resp, nil
}
