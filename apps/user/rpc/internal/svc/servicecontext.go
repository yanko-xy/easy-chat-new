package svc

import (
	"rpc/internal/config"

	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModels models.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		UserModels: models.NewUserModel(sqlConn, c.Cache),
	}
}
