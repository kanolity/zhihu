package code

import "go_code/zhihu/pkg/xcode"

var (
	RegisterMobileInvalid   = xcode.New(10001, "无效的手机号")
	RegisterPasswordEmpty   = xcode.New(10002, "密码不能为空")
	VerificationCodeEmpty   = xcode.New(10003, "验证码不能为空")
	UsernameHasRegistered   = xcode.New(10004, "用户名已注册")
	MobileHasRegistered     = xcode.New(10005, "手机号已注册")
	LoginUsernameEmpty      = xcode.New(10006, "用户名为空")
	LoginPasswordEmpty      = xcode.New(10007, "密码为空")
	VerificationCodeExpired = xcode.New(10008, "验证码已过期")
	VerificationCodeFailed  = xcode.New(10009, "验证码错误")
	LoginPasswordFailed     = xcode.New(10010, "密码错误")
	LoginUserUnexisted      = xcode.New(10011, "用户不存在")
)
