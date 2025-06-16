package code

import "go_code/zhihu/pkg/xcode"

var (
	RegisterNameEmpty   = xcode.New(20001, "用户名不能为空")
	RegisterNameExisted = xcode.New(20002, "用户名已存在")
)
