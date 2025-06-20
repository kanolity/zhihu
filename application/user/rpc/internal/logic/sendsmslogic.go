package logic

import (
	"context"
	"go_code/zhihu/application/user/rpc/internal/code"
	"go_code/zhihu/application/user/rpc/internal/svc"
	"go_code/zhihu/application/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSmsLogic) SendSms(in *user.SendSmsRequest) (*user.SendSmsResponse, error) {
	// 校验手机号
	if len(in.Mobile) != 11 {
		return nil, code.VerficationMobileInvalid
	}

	// 记录日志
	logx.WithContext(l.ctx).Infof("发送验证码, Mobile: %s, Code: %s", in.Mobile, in.Code)

	// 调用短信服务
	//err := l.svcCtx.SmsClient.Send(in.Mobile, in.Code)
	//if err != nil {
	//	logx.WithContext(l.ctx).Errorf("短信发送失败，Mobile: %s, err: %v", in.Mobile, err)
	//	return nil, status.Error(codes.Internal, "短信发送失败")
	//}

	return &user.SendSmsResponse{}, nil
}
