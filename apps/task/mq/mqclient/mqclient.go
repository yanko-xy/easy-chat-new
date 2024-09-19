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

type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addrs []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addrs, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(context.Background(), string(body))
}
