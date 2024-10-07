/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type LimitMiddleware struct {
	redisCfg redis.RedisConf
	*limit.TokenLimiter
	*limit.PeriodLimit
}

func NewLimitMiddleware(cfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{
		redisCfg: cfg,
	}
}

func (m *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	m.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(m.redisCfg), "REDIS_TOKEN_LIMIT_KEY")
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if m.TokenLimiter.AllowCtx(request.Context()) {
				next(writer, request)
				return
			}
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

func (m *LimitMiddleware) PeriodLimitHandle(seconds, quota int) rest.Middleware {
	m.PeriodLimit = limit.NewPeriodLimit(seconds, quota, redis.MustNewRedis(m.redisCfg), "REDIS_PERIOD_LIMIT_KEY", limit.Align())
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			val, _ := m.PeriodLimit.TakeCtx(request.Context(), "first")
			if val == limit.Allowed {
				next(writer, request)
				return
			}
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}
