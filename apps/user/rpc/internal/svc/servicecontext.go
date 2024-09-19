package svc

import (
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/config"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"

	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	*redis.Redis
	UserModels models.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UserModels: models.NewUserModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 9999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	// 写入到redis
	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
