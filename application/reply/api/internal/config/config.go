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
	ReplyRPC   zrpc.RpcClientConf
	MessageRPC zrpc.RpcClientConf
	ArticleRPC zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf
}
