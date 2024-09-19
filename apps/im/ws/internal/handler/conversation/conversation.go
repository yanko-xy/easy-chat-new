/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package conversation

import (
	"github.com/mitchellh/mapstructure"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/apps/task/mq/mq"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			srv.Errorf("websocket conn mapstructure decode err %v, msg %v", err, msg.Data)
			return
		}

		switch data.ChatType {
		case constants.SingleChatTpye:
			err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				SendTime:       time.Now().UnixNano(),
				MType:          data.Msg.MType,
				Content:        data.Msg.Content,
			})
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				srv.Errorf("mq push  err %v", err)
				return
			}
		}
	}
}
