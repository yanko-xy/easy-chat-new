/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
)

type Message struct {
	FrameType FrameType   `json:"frameType"`
	Method    string      `json:"method"`
	FormId    string      `json:"formId"`
	Data      interface{} `json:"data"`
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}
