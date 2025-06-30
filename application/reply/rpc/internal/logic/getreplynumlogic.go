package logic

import (
	"context"
	"go_code/zhihu/application/reply/rpc/internal/svc"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReplyNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyNumLogic {
	return &GetReplyNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReplyNumLogic) GetReplyNum(in *reply.GetReplyNumReq) (*reply.GetReplyNumResp, error) {
	replyNum, err := l.svcCtx.ReplyCountModel.FindByArticleId(l.ctx, in.ArticleId)
	if err != nil {
		return nil, err
	}
	return &reply.GetReplyNumResp{
		ReplyNum: replyNum,
	}, nil
}
