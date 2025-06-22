package code

import "go_code/zhihu/pkg/xcode"

var (
	TitleCantEmpty            = xcode.New(30002, "标题不能为空")
	ArticleContentTooFewWords = xcode.New(30003, "文章内容字数太少")
	GetArticleDetailFailed    = xcode.New(30004, "获取文章详情失败")
	GetUserInfoFailed         = xcode.New(30005, "获取作者信息失败")
)
