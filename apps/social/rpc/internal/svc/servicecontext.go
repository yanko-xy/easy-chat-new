package svc

import (
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/config"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	socialmodels.FriendModel
	socialmodels.FriendRequestModel
	socialmodels.GroupModel
	socialmodels.GroupMemberModel
	socialmodels.GroupRequestModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		FriendModel:        socialmodels.NewFriendModel(sqlConn, c.Cache),
		FriendRequestModel: socialmodels.NewFriendRequestModel(sqlConn, c.Cache),
		GroupModel:         socialmodels.NewGroupModel(sqlConn, c.Cache),
		GroupMemberModel:   socialmodels.NewGroupMemberModel(sqlConn, c.Cache),
		GroupRequestModel:  socialmodels.NewGroupRequestModel(sqlConn, c.Cache),
	}
}
