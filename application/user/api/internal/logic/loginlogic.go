package logic

import (
	"context"
	"go_code/zhihu/application/user/api/internal/code"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/pkg/encrypt"
	"go_code/zhihu/pkg/jwt"
	"go_code/zhihu/pkg/xcode"
	"strings"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) == 0 {
		return nil, code.LoginUsernameEmpty
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.LoginPasswordEmpty
	}
	u, err := l.checkPassword(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	if u == nil || u.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": u.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		UserId: u.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

func (l *LoginLogic) checkPassword(username string, password string) (*user.FindByUsernameResponse, error) {
	password = encrypt.EncPassword(password)
	user1, err := l.svcCtx.UserRpc.FindByUsername(l.ctx, &user.FindByUsernameRequest{Username: username})
	if err != nil {
		logx.Errorf("FindByUsername %s error: %v", username, err)
		return nil, err
	}
	if user1 == nil {
		return nil, code.LoginUserUnexisted
	}
	if user1.Password != password {
		return nil, code.LoginPasswordFailed
	}
	return user1, nil
}
