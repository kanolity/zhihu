package types

const (
	SortPublishTime = iota
	SortLikeCount
)

const (
	DefaultPageSize = 10
	DefaultLimit    = 100

	DefaultSortLikeCursor = 1 << 30
)

const (
	// ArticleStatusPending 待审核
	ArticleStatusPending = iota
	// ArticleStatusNotPass 审核不通过
	ArticleStatusNotPass
	// ArticleStatusVisible 可见
	ArticleStatusVisible
	// ArticleStatusUserDelete 用户删除
	ArticleStatusUserDelete
)
