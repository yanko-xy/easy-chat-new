package main

import (
	"flag"
	"fmt"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/config"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/server"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"github.com/yanko-xy/easy-chat/pkg/configserver"
	"sync"

	"github.com/yanko-xy/easy-chat/pkg/interceptor/rpcserver"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

var grpcServer *grpc.Server
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigserver(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "120.26.209.19:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "user",
		Configs:        "user-rpc.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		grpcServer.GracefulStop()
		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()
			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}
	// 设置日志配置
	//var lc logx.LogConf
	//conf.MustLoad(*logConfigFile, &lc)
	//logx.MustSetup(lc)
	//logx.AddWriter(logx.NewWriter(os.Stdout))

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()
		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
		grpcServer = srv

		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	if err := ctx.SetRootToken(); err != nil {
		fmt.Errorf("%v", err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
