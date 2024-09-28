/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package handler

import (
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler/conversation"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler/push"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/handler/user"
	"github.com/yanko-xy/easy-chat/apps/im/ws/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svcCtx *svc.ServiceContext) {
	srv.AddRoutes([]*websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svcCtx),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svcCtx),
		},
		{
			Method:  "conversation.markRead",
			Handler: conversation.MarkRead(svcCtx),
		},

		{
			Method:  "push",
			Handler: push.Push(svcCtx),
		},
	})
}
