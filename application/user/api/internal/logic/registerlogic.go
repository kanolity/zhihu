package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go_code/zhihu/application/user/api/internal/code"
	"go_code/zhihu/application/user/rpc/types/user"
	"go_code/zhihu/pkg/encrypt"
	"go_code/zhihu/pkg/jwt"
	"strings"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixActivation = "biz#activation#%s"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	//标准化用户输入
	req.Username = strings.TrimSpace(req.Username)
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) != 11 {
		return nil, code.RegisterMobileInvalid
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.RegisterPasswordEmpty
	} else {
		//加密密码
		req.Password = encrypt.EncPassword(req.Password)
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	//校验验证码
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Infof("check verification code err: %v", err)
		return nil, err
	}

	//加密手机号
	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Infof("encode mobile[%s] err: %v", req.Mobile, err)
		return nil, err
	}

	//查询手机号及名字是否已被注册
	u1, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: req.Mobile})
	if err != nil {
		logx.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	if u1 != nil && u1.UserId > 0 {
		return nil, code.MobileHasRegistered
	}

	u2, err := l.svcCtx.UserRpc.FindByUsername(l.ctx, &user.FindByUsernameRequest{Username: req.Username})
	if err != nil {
		logx.Errorf("FindByUsername error: %v", err)
		return nil, err
	}
	if u2 != nil && u2.UserId > 0 {
		return nil, code.UsernameHasRegistered
	}

	//注册及生成token
	regMsg, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Username,
		Mobile:   mobile,
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": regMsg.UserId,
		},
	})
	if err != nil {
		logx.Errorf("Build token error: %v", err)
		return nil, err
	}
	_ = delActivationCache(req.Mobile, req.VerificationCode, l.svcCtx.BizRedis)
	return &types.RegisterResponse{
		UserId: regMsg.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

func checkVerificationCode(rds *redis.Redis, mobile, code string) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}
	return nil
}
