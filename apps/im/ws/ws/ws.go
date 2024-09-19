/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package ws

import "github.com/yanko-xy/easy-chat/pkg/constants"

type (
	Msg struct {
		MType   constants.MType `mapstructure:"mType" json:"mType"`
		Content string          `mapstructure:"content" json:"content"`
	}

	Chat struct {
		ConversationId string             `mapstructure:"conversationId" json:"conversationId"`
		ChatType       constants.ChatType `mapstructure:"chatType" json:"chatType"`
		SendId         string             `mapstructure:"sendId" json:"sendId"`
		RecvId         string             `mapstructure:"recvId" json:"recvId"`
		Msg            Msg                `mapstructure:"msg" json:"msg"`
		SendTime       int64              `mapstructure:"sendTime" json:"sendTime"`
	}
)
