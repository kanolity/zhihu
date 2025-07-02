package logic

import (
	"context"
	"fmt"
	"go_code/zhihu/application/article/rpc/types/article"
	"go_code/zhihu/application/message/rpc/types/message"
	"go_code/zhihu/application/reply/rpc/types/reply"
	"go_code/zhihu/application/user/rpc/types/user"

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
		logx.Errorf("Post Reply  error: %v", err)
		return nil, err
	}
	if req.BizId == "article" {
		_, err = l.svcCtx.ArticleRpc.ArticleReplyIncrease(l.ctx, &article.ArticleReplyIncreaseRequest{
			ArticleId: req.TargetId,
		})
		if err != nil {
			logx.Errorf("increase article[%v] comment count  error: %v", req.TargetId, err)
		}
	}
	replyedUsername, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: req.ReplyUserId,
	})
	if err != nil {
		logx.Errorf("find user by id error: %v", err)
	}
	_, err = l.svcCtx.MessageRpc.SendMessage(l.ctx, &message.SendMessageRequest{
		Type:       2,
		BizId:      "reply",
		TargetId:   req.TargetId,
		ReceiverId: req.BeReplyUserId,
		Title:      "收到回复",
		Content:    fmt.Sprintf("用户%s回复:%s", replyedUsername.Username, req.Content),
	})
	if err != nil {
		logx.Errorf("post reply err:%v", err)
	}
	return &types.PostReplyResp{Id: response.Id}, nil
}
