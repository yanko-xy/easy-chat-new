/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package mq

import "github.com/yanko-xy/easy-chat/pkg/constants"

type MsgChatTransfer struct {
	ConversationId string             `json:"conversationId"`
	ChatType       constants.ChatType `json:"chatType"`
	SendId         string             `json:"sendId"`
	RecvId         string             `json:"recvId"`
	SendTime       int64              `json:"sendTime"`
	MType          constants.MType    `json:"mType"`
	Content        string             `json:"content"`
}
