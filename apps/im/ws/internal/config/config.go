/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf

	ListenOn string

	JwtAuth struct {
		AccessSecret string
	}

	Mongo struct {
		Url string
		Db  string
	}

	Redisx redis.RedisConf

	MsgChatTransfer struct {
		Topic string
		Addrs []string
	}

	MsgMarkRead struct {
		Topic string
		Addrs []string
	}
}
