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
	"github.com/yanko-xy/easy-chat/pkg/configserver"
	"github.com/zeromicro/go-zero/core/service"
	"sync"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigserver(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "120.26.209.19:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
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
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}

	fmt.Println("start mqueue at", c.ListenOn, "...")
	serviceGroup.Start()
}
