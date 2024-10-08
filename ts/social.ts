import webapi from "./gocliRequest"
import * as components from "./socialComponents"
export * from "./socialComponents"

/**
 * @description "好友申请"
 * @param req
 */
export function friendPutIn(req: components.FriendPutInReq) {
	return webapi.post<components.FriendPutInResp>(`/v1/social/friend/putIn`, req)
}

/**
 * @description "好友申请处理"
 * @param req
 */
export function friendPutInHandle(req: components.FriendPutInHandleReq) {
	return webapi.put<components.FriendPutInHandleResp>(`/v1/social/friend/putIn`, req)
}

/**
 * @description "好友申请列表"
 */
export function friendPutInList() {
	return webapi.get<components.FriendPutInListResp>(`/v1/social/friend/putIns`)
}

/**
 * @description "好友列表"
 */
export function friendList() {
	return webapi.get<components.FriendListResp>(`/v1/social/friends`)
}

/**
 * @description "好友在线情况"
 */
export function friendsOnline() {
	return webapi.get<components.FriendsOnlineResp>(`/v1/social/friends/online`)
}

/**
 * @description "创群"
 * @param req
 */
export function createGroup(req: components.GroupCreateReq) {
	return webapi.post<components.GroupCreateResp>(`/v1/social/group`, req)
}

/**
 * @description "申请进群"
 * @param req
 */
export function groupPutIn(req: components.GroupPutInReq) {
	return webapi.post<components.GroupPutInResp>(`/v1/social/group/putIn`, req)
}

/**
 * @description "申请进群处理"
 * @param req
 */
export function groupPutInHandle(req: components.GroupPutInHandleReq) {
	return webapi.put<components.GroupPutInHandleResp>(`/v1/social/group/putIn`, req)
}

/**
 * @description "申请进群列表"
 * @param params
 */
export function groupPutInList(params: components.GroupPutInListReqParams) {
	return webapi.get<components.GroupPutInListResp>(`/v1/social/group/putIns`, params)
}

/**
 * @description "成员列表列表"
 * @param params
 */
export function groupUserList(params: components.GroupUserListReqParams) {
	return webapi.get<components.GroupUserListResp>(`/v1/social/group/users`, params)
}

/**
 * @description "群在线用户"
 * @param params
 */
export function groupUserOnline(params: components.GroupUserOnlineReqParams) {
	return webapi.get<components.GroupUserOnlineResp>(`/v1/social/group/users/online`, params)
}

/**
 * @description "用户申群列表"
 */
export function groupList() {
	return webapi.get<components.GroupListResp>(`/v1/social/groups`)
}
