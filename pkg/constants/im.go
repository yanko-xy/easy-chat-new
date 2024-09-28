/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package constants

type MType int

const (
	TextMType MType = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

type ContentType int

const (
	ContentChatMsg ContentType = iota
	ContentMakeRead
)
