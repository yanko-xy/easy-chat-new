/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import "time"

type FrameType uint8

const (
	FrameData      FrameType = 0x0
	FramePing      FrameType = 0x1
	FrameAck       FrameType = 0x2
	FrameNoAck     FrameType = 0x3
	FrameTranspond FrameType = 0x6
	FrameErr       FrameType = 0x9
)

type Message struct {
	FrameType    FrameType `json:"frameType"`
	Id           string    `json:"id"`
	TranspondUid string    `json:"transpondUid"`
	AckSeq       int       `json:"ackSeq"`
	ackTime      time.Time
	errCount     int
	Method       string      `json:"method"`
	FormId       string      `json:"formId"`
	Data         interface{} `json:"data"` // map[string]interface{}
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
