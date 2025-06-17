package code

import "go_code/zhihu/pkg/xcode"

var (
	RegisterNameEmpty    = xcode.New(20001, "用户名不能为空")
	FindUserFailed       = xcode.New(20002, "查询用户失败")
	ChangePasswordWrong  = xcode.New(20003, "密码错误")
	ChangePasswordFailed = xcode.New(20004, "更新密码错误")
	ChangeAvatarFailed   = xcode.New(20005, "头像更新失败")
)
