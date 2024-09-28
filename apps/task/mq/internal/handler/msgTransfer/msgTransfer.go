/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package msgTransfer

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/socialclient"
	"github.com/yanko-xy/easy-chat/apps/task/mq/internal/svc"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseTransfer(svc *svc.ServiceContext) *baseTransfer {
	return &baseTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *baseTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = m.group(ctx, data)
	case constants.SingleChatType:
		err = m.single(ctx, data)
	}
	return err
}

func (m *baseTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *baseTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群的用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})

	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))
	for _, member := range users.List {
		// 跳过发送人
		if member.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, member.UserId)
	}

	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
