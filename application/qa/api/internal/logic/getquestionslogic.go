package logic

import (
	"context"
	"go_code/zhihu/application/qa/rpc/types/qa"
	"go_code/zhihu/application/user/rpc/types/user"

	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionsLogic {
	return &GetQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionsLogic) GetQuestions(req *types.GetQuestionsReq) (resp *types.GetQuestionsResp, err error) {
	response, err := l.svcCtx.QaRpc.GetQuestions(l.ctx, &qa.GetQuestionsRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		Limit:      req.Limit,
	})
	if err != nil {
		return nil, err
	}

	questions := response.Questions

	// 提取所有用户 ID 去重
	userIdSet := make(map[int64]struct{})
	for _, q := range questions {
		userIdSet[q.UserId] = struct{}{}
	}

	userIds := make([]int64, 0, len(userIdSet))
	for id := range userIdSet {
		userIds = append(userIds, id)
	}

	// 批量查询用户信息
	userResp, err := l.svcCtx.UserRpc.BatchGetUsers(l.ctx, &user.BatchGetUsersRequest{
		UserIds: userIds,
	})
	if err != nil {
		logx.Errorf("BatchGetUsers failed: %v", err)
		return nil, err
	}

	// 构建 userId → user 映射
	userMap := make(map[int64]*user.UserInfo)
	for _, u := range userResp.Users {
		userMap[u.Id] = u
	}

	// 组装最终返回值
	result := make([]types.GetQuestion, 0, len(questions))
	for _, q := range questions {
		u := userMap[q.UserId]
		result = append(result, types.GetQuestion{
			Id:          q.Id,
			Title:       q.Title,
			Username:    u.GetUsername(),
			Avatar:      u.GetAvatar(),
			Description: q.Description,
			IsResolved:  q.IsResolved,
			CreateTime:  q.CreateTime,
		})
	}

	return &types.GetQuestionsResp{
		Questions: result,
		HasMore:   response.HasMore,
	}, nil
}
