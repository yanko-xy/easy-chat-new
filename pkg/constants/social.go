/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package constants

type HandleResult int

// 处理结果 1. 未处理，2. 处理， 3. 拒绝  4. 取消
const (
	NoHandleResult HandleResult = iota + 1
	PassHandleResutl
	RefuseHandleResult
	CancelHandleResult
)
