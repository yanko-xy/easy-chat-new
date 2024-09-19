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
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			srv.Errorf("websocket conn mapstructure decode err %v, msg %v", err, msg.Data)
			return
		}

		// 发送的目标
		rconn := srv.GetConn(data.RecvId)
		if rconn == nil {
			// todo: 目标离线
			return
		}

		srv.Infof("push msg %v", data)

		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         data.SendId,
			RecvId:         data.RecvId,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
		}), rconn)
	}
}
