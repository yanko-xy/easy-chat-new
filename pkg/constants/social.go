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

// 群等级
// 1. 创建者
// 2. 管理者
// 3. 普通
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1
	ManagerGroupRoleLevel
	AtLargeGroupRoleLevel
)

// 进群申请方式
// 1. 邀请  2. 申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)
