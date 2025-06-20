package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go_code/zhihu/application/user/api/internal/code"
	"go_code/zhihu/application/user/rpc/userclient"
	"go_code/zhihu/pkg/utils"
	"strconv"
	"time"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	//获取今日验证码获取次数
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("get verification count mobile[%s] err:%v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		return nil, code.VerificateTooMany
	}

	// 获取验证码，若无则创建
	vCode, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("get activation cache mobile[%s] err:%v", req.Mobile, err)
	}
	if len(vCode) == 0 {
		vCode = utils.RandomNumber(6)
	}

	//将验证码存入缓存
	err = saveActivationCache(req.Mobile, vCode, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("save activation cache mobile[%s] err:%v", req.Mobile, err)
		return nil, err
	}

	//发送短信
	_, err = l.svcCtx.UserRpc.SendSms(l.ctx, &userclient.SendSmsRequest{
		Mobile: req.Mobile,
		Code:   vCode,
	})
	if err != nil {
		logx.Errorf("send sms mobile [%s] err:%v", req.Mobile, err)
		return nil, err
	}

	//增加今日验证码获取次数
	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incr verification count mobile[%s] err:%v", req.Mobile, err)
	}
	return &types.VerificationResponse{}, nil
}

// getVerificationCount 获取今日验证码获取次数
func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// incrVerificationCount 增加今日验证码获取次数
func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expireat(key, utils.EndOfDay(time.Now()).Unix())
}

// getActivationCache 获取验证码
func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

// saveActivationCache 将验证码存入缓存
func saveActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	fmt.Println(key)
	return rds.Setex(key, code, expireActivation)
}

// delActivationCache 删除验证码
func delActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
