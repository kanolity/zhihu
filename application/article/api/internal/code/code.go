package code

import "go_code/zhihu/pkg/xcode"

var (
	TitleCantEmpty            = xcode.New(30002, "标题不能为空")
	ArticleContentTooFewWords = xcode.New(30003, "文章内容字数太少")
)
