package main

import (
	"flag"
	"fmt"

	"go_code/zhihu/application/qa/rpc/internal/config"
	"go_code/zhihu/application/qa/rpc/internal/server"
	"go_code/zhihu/application/qa/rpc/internal/svc"
	"go_code/zhihu/application/qa/rpc/types/qa"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/qa.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		qa.RegisterQaServiceServer(grpcServer, server.NewQaServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
