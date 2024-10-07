package svc

import (
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/config"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/socialclient"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/userclient"
	"github.com/yanko-xy/easy-chat/pkg/interceptor"
	"github.com/yanko-xy/easy-chat/pkg/middleware"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

type ServiceContext struct {
	Config config.Config

	socialclient.Social
	userclient.User

	IdempotenceMiddleware rest.Middleware
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotenceClient),
		)),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),

		IdempotenceMiddleware: middleware.NewIdempotenceMiddle().Handle,
		Redis:                 redis.MustNewRedis(c.Redisx),
	}
}
