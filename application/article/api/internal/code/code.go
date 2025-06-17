package code

import "go_code/zhihu/pkg/xcode"

var (
	UploadCoverFailed         = xcode.New(30001, "保存封面失败")
	TitleCantEmpty            = xcode.New(30002, "标题不能为空")
	ArticleContentTooFewWords = xcode.New(30003, "文章内容字数太少")
	CoverCantEmpty            = xcode.New(30004, "封面不能为空")
)
