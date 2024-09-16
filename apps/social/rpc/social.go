package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/config"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/server"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")
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
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
