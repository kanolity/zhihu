package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
	}
	ArticleRPC zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf
	TagRPC     zrpc.RpcClientConf
	MessageRPC zrpc.RpcClientConf
	Es         struct {
		Addresses []string
		Username  string
		Password  string
	}
}
