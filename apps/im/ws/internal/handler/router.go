/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package handler

import (
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler/conversation"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler/user"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]*websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
	})
}
