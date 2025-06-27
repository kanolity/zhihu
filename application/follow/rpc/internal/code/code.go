package code

import "go_code/zhihu/pkg/xcode"

var (
	FollowUserIdEmpty   = xcode.New(80001, "关注用户id为空")
	FollowedUserIdEmpty = xcode.New(80002, "被关注用户id为空")
	CannotFollowSelf    = xcode.New(80003, "不能关注自己")
	UserIdEmpty         = xcode.New(80004, "用户id为空")
)
