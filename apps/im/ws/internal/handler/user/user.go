/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package user

import (
	"github.com/gorilla/websocket"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	websocketx "github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(u[0], uids), conn)
		if err != nil {
			srv.Error("err ", err)
		}
	}
}
