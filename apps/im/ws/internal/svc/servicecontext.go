/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package svc

import (
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel: immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
