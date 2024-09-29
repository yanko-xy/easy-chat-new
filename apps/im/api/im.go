package main

import (
	"flag"
	"fmt"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/config"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/handler"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/svc"
	"github.com/yanko-xy/easy-chat/pkg/configserver"
	"github.com/yanko-xy/easy-chat/pkg/resultx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"sync"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigserver(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "120.26.209.19:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-api.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		proc.WrapUp()
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

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()
		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
