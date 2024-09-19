/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package main

import (
	"flag"
	"fmt"
	"github.com/yanko-xy/easy-chat/apps/task/mq/internal/config"
	"github.com/yanko-xy/easy-chat/apps/task/mq/internal/handler"
	"github.com/yanko-xy/easy-chat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"os"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")
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

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}

	fmt.Println("start mqueue at", c.ListenOn, "...")
	serviceGroup.Start()
}
