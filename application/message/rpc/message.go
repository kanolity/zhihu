package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go_code/zhihu/pkg/xcode"

	"go_code/zhihu/application/message/rpc/internal/config"
	"go_code/zhihu/application/message/rpc/internal/server"
	"go_code/zhihu/application/message/rpc/internal/svc"
	"go_code/zhihu/application/message/rpc/types/message"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/message.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		message.RegisterMessageServiceServer(grpcServer, server.NewMessageServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	httpx.SetErrorHandler(xcode.ErrHandler)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
