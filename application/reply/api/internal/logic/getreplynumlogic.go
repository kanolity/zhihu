package logic

import (
	"context"
	"go_code/zhihu/application/reply/rpc/types/reply"

	"go_code/zhihu/application/reply/api/internal/svc"
	"go_code/zhihu/application/reply/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReplyNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyNumLogic {
	return &GetReplyNumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReplyNumLogic) GetReplyNum(req *types.GetReplytNumReq) (resp *types.GetReplytNumResp, err error) {
	response, err := l.svcCtx.ReplyRpc.GetReplyNum(l.ctx, &reply.GetReplyNumReq{ArticleId: req.ArticleId})
	if err != nil {
		return nil, err
	}
	return &types.GetReplytNumResp{
		CommentNum: response.ReplyNum,
	}, nil
}
