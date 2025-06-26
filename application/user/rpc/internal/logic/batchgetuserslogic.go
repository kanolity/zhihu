package logic

import (
	"context"

	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUsersLogic) BatchGetUsers(in *user.BatchGetUsersRequest) (*user.BatchGetUsersResponse, error) {
	if len(in.UserIds) == 0 {
		return &user.BatchGetUsersResponse{Users: []*user.UserInfo{}}, nil
	}

	users, err := l.svcCtx.UserModel.FindUsersByIds(l.ctx, in.UserIds)
	if err != nil {
		return nil, err
	}

	resp := &user.BatchGetUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &user.UserInfo{
			Id:       u.Id,
			Username: u.Username,
			Avatar:   u.Avatar,
		})
	}
	return resp, nil
}
