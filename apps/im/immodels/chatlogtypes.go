/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package immodels

import (
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChatLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string             `bson:"conversationId"`
	SendId         string             `bson:"sendId"`
	RecvId         string             `bson:"recvId"`
	MsgFrom        int                `bson:"msgFrom"`
	ChatType       constants.ChatType `bson:"chatType"`
	MsgType        constants.MType    `bson:"msgType"`
	MsgContent     string             `bson:"msgContent"`
	SendTime       int64              `bson:"sendTime"`
	Status         int                `bson:"status"`

	// TODO: Fill other fields
	UpdateAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreateAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}
