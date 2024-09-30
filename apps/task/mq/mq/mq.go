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
	RecvIds        []string           `json:"recvIds"`
	SendTime       int64              `json:"sendTime"`
	MsgId          string             `json:"msgId"`
	MType          constants.MType    `json:"mType"`
	Content        string             `json:"content"`
}

type MsgMarkRead struct {
	// 消息类型：1. 私聊、2. 群聊
	ChatType constants.ChatType `json:"chatType,omitempty"`
	// 会话id
	ConversationId string `json:"conversationId,omitempty"`
	// 发送者
	SendId string `json:"sendId,omitempty"`
	// 接收者
	RecvId string `json:"recvId,omitempty"`
	// 已读消息集合
	MsgIds []string `json:"msgIds,omitempty"`
}
