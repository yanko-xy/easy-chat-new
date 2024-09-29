/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package main

import (
	"flag"
	"fmt"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/config"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
	"github.com/yanko-xy/easy-chat/pkg/configserver"
	"sync"
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
		Configs:        "im-ws.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

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
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthorization(handler.NewJwtAuth(ctx)),
		//websocket.WithServerAck(websocket.RigorAck),
		// 心跳检测
		//websocket.WithServerMaxConnectionIdle(10*time.Second),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at", c.ListenOn, "...")
	srv.Start()
}
