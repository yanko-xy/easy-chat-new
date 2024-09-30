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
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"net/http"
	"sync"
	"time"
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
	// 设置服务认证的token
	token, err := ctxdata.GetJwtToken(c.JwtAuth.AccessSecret, time.Now().Unix(), 3153600000, fmt.Sprintf("ws:%s", time.Now().Unix()))
	if err != nil {
		panic(err)
	}

	opts := []websocket.ServerOptions{
		websocket.WithServerAuthorization(handler.NewJwtAuth(ctx)),
		websocket.WithServerDiscover(websocket.NewRedisDiscover(http.Header{
			"Authorization": []string{token},
		}, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	}

	srv := websocket.NewServer(c.ListenOn, opts...)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at", c.ListenOn, "...")
	srv.Start()
}
