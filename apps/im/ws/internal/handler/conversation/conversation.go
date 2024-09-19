/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package conversation

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/logic"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		switch data.ChatType {
		case constants.SingleChatTpye:
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}

			srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				Msg:            data.Msg,
				SendTime:       time.Now().UnixMilli(),
			}), data.RecvId)
		}
	}
}
