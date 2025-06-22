package code

import "go_code/zhihu/pkg/xcode"

var (
	GetIsThumbupFailed = xcode.New(50001, "查询是否点赞失败")
	ThumbupFailed      = xcode.New(50002, "点赞失败")
)
