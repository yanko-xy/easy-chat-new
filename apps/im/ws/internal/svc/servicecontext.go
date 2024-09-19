/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package svc

import (
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/config"
	"github.com/yanko-xy/easy-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
