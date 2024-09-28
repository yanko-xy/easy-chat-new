/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package push

import (
	"github.com/mitchellh/mapstructure"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/pkg/constants"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var (
			data ws.Push
			err  error
		)
		if err = mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			srv.Errorf("websocket conn mapstructure decode err %v, msg %v", err, msg.Data)
			return
		}

		srv.Infof("push req %v", data)
		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			group(srv, &data)
		}

	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		srv.Infof("push conn off line %v", recvId)
		// todo: 目标离线
		return nil
	}

	srv.Infof("push uid %v", recvId)
	err := srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			ReadRecords: data.ReadRecords,
			MType:       data.MType,
			Content:     data.Content,
		},
	}), rconn)

	if err != nil {
		srv.Errorf("push err %v, recvId %v", err, recvId)
	}

	return nil
}

// 基于并发发送
func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}

	return nil
}
