package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionsLogic {
	return &GetQuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionsLogic) GetQuestions(in *qa.GetQuestionsRequest) (*qa.GetQuestionsResponse, error) {
	// 设定 pageSize +1 获取下一页判断
	limit := in.Limit + 1

	// 使用 cursor 和 question_id 实现双重排序分页
	// 假设是按 create_time + id 倒序（示意）
	questions, err := l.svcCtx.QuestionModel.ListQuestionsByCursor(l.ctx, in.Cursor, in.QuestionId, limit)
	if err != nil {
		logx.Errorf("GetQuestions failed: %v", err)
		return nil, status.Error(codes.Internal, "数据库查询失败")
	}

	hasMore := false
	if int64(len(questions)) == limit {
		hasMore = true
		questions = questions[:len(questions)-1] // 截取实际数量
	}

	list := make([]*qa.Question, 0, len(questions))
	for _, q := range questions {
		list = append(list, &qa.Question{
			Id:          int64(q.Id),
			UserId:      int64(q.UserId),
			Title:       q.Title,
			Description: q.Description,
			IsResolved:  q.IsResolved,
			CreateTime:  q.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &qa.GetQuestionsResponse{
		Questions: list,
		HasMore:   hasMore,
	}, nil
}
