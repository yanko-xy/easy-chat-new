/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package svc

import (
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/config"
	"github.com/yanko-xy/easy-chat/apps/task/mq/mqclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgMarkRead.Addrs, c.MsgMarkRead.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		Redis:                 redis.MustNewRedis(c.Redisx),
	}
}
