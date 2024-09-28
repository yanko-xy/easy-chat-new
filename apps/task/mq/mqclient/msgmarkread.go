/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package mqclient

import (
	"context"
	"encoding/json"
	"github.com/yanko-xy/easy-chat/apps/task/mq/mq"
	"github.com/zeromicro/go-queue/kq"
)

type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

type msgReadTransferCilent struct {
	pusher *kq.Pusher
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) *msgReadTransferCilent {
	return &msgReadTransferCilent{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgReadTransferCilent) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(context.Background(), string(body))
}
