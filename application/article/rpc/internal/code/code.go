package code

import "go_code/zhihu/pkg/xcode"

var (
	UserIdInvalid     = xcode.New(40001, "用户id无效")
	SortTypeInvalid   = xcode.New(40002, "排序类型无效")
	FindArticleFailed = xcode.New(40003, "获取文章信息失败")
	ArticleIdInvalid  = xcode.New(40005, "文章ID无效")
)
