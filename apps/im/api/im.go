package main

import (
	"flag"
	"fmt"
	"github.com/yanko-xy/easy-chat/pkg/resultx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"os"

	"github.com/yanko-xy/easy-chat/apps/im/api/internal/config"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/handler"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
