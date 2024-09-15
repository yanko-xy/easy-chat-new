package main

import (
	"flag"
	"fmt"
	"os"

	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/internal/svc"
	"rpc/user"

	"github.com/yanko-xy/easy-chat/pkg/interceptor/rpcserver"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")
var logConfigFile = flag.String("log", "../../etc/log.yaml", "log config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 设置日志配置
	var lc logx.LogConf
	conf.MustLoad(*logConfigFile, &lc)
	logx.MustSetup(lc)
	logx.AddWriter(logx.NewWriter(os.Stdout))

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
