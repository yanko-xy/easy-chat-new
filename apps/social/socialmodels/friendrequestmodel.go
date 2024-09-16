package socialmodels

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestModel = (*customFriendRequestModel)(nil)

type (
	// FriendRequestModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestModel.
	FriendRequestModel interface {
		friendRequestModel
	}

	customFriendRequestModel struct {
		*defaultFriendRequestModel
	}
)

// NewFriendRequestModel returns a model for the database table.
func NewFriendRequestModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestModel {
	return &customFriendRequestModel{
		defaultFriendRequestModel: newFriendRequestModel(conn, c, opts...),
	}
}
