// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3

package types

type GetRepliesReq struct {
	BizId    string `form:"biz_id"`
	TargetId int64  `form:"target_id"`
	Cursor   int64  `form:"cursor"`
	Limit    int64  `form:"limit"`
}

type GetRepliesResp struct {
	Replies []Reply `json:"replies"`
	HasMore bool    `json:"has_more"`
}

type GetReplytNumReq struct {
	ArticleId int64 `form:"article_id"`
}

type GetReplytNumResp struct {
	CommentNum int64 `json:"comment_num"`
}

type PostReplyReq struct {
	BizId         string `json:"biz_id"`
	TargetId      int64  `json:"target_id"`
	ReplyUserId   int64  `json:"reply_user_id"`
	BeReplyUserId int64  `json:"be_reply_user_id"`
	ParentId      int64  `json:"parent_id"`
	Content       string `json:"content"`
}

type PostReplyResp struct {
	Id int64 `json:"id"`
}

type Reply struct {
	Id            int64  `json:"id"`
	BizId         string `json:"biz_id"`
	TargetId      int64  `json:"target_id"`
	ReplyUserId   int64  `json:"reply_user_id"`
	BeReplyUserId int64  `json:"be_reply_user_id"`
	ParentId      int64  `json:"parent_id"`
	Content       string `json:"content"`
	LikeNum       int64  `json:"like_num"`
	IsDeleted     bool   `json:"is_deleted"`
	CreateTime    string `json:"create_time"`
}
