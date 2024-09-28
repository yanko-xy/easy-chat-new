/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package ws

import "github.com/yanko-xy/easy-chat/pkg/constants"

type (
	Msg struct {
		MsgId       string            `mapstructure:"msgId" json:"msgId"`
		ReadRecords map[string]string `mapstructure:"readRecords" json:"readRecords"`
		MType       constants.MType   `mapstructure:"mType" json:"mType"`
		Content     string            `mapstructure:"content" json:"content"`
	}

	Chat struct {
		ConversationId string             `mapstructure:"conversationId" json:"conversationId"`
		ChatType       constants.ChatType `mapstructure:"chatType" json:"chatType"`
		SendId         string             `mapstructure:"sendId" json:"sendId"`
		RecvId         string             `mapstructure:"recvId" json:"recvId"`
		SendTime       int64              `mapstructure:"sendTime" json:"sendTime"`
		Msg            `mapstructure:"msg"`
	}

	Push struct {
		ConversationId string             `mapstructure:"conversationId" json:"conversationId"`
		ChatType       constants.ChatType `mapstructure:"chatType" json:"chatType"`
		SendId         string             `mapstructure:"sendId" json:"sendId"`
		RecvId         string             `mapstructure:"recvId" json:"recvId"`
		RecvIds        []string           `mapstructure:"recvIds" json:"recvIds"`
		SendTime       int64              `mapstructure:"sendTime" json:"sendTime"`

		MsgId       string                `mapstructure:"msgId" json:"msgId"`
		ReadRecords map[string]string     `mapstructure:"readRecords" json:"readRecords"`
		ContentType constants.ContentType `mapstructure:"contentType" json:"contentType"`

		MType   constants.MType `mapstructure:"mType" json:"mType"`
		Content string          `mapstructure:"content" json:"content"`
	}

	MaskRead struct {
		ChatType       constants.ChatType `mapstructure:"chatType" json:"chatType"`
		RecvId         string             `mapstructure:"recvId" json:"recvId"`
		ConversationId string             `mapstructure:"conversationId" json:"conversationId"`
		MsgIds         []string           `mapstructure:"msgIds" json:"msgIds"`
	}
)
