package group

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/userclient"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 成员列表列表
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
		GroupId: req.GroupId,
	})

	// 获取用户信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return nil, err
	}

	userRecord := make(map[string]*userclient.UserEntity, len(userList.Users))
	for i, v := range userList.Users {
		userRecord[v.Id] = userList.Users[i]
	}

	respList := make([]*types.GroupMember, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		member := &types.GroupMember{
			Id:          int64(v.Id),
			GroupId:     v.GroupId,
			UserId:      v.UserId,
			RoleLevel:   int(v.RoleLevel),
			InviterUid:  v.InviterUid,
			OperatorUid: v.OperatorUid,
		}
		if u, ok := userRecord[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}

		respList = append(respList, member)
	}
	return &types.GroupUserListResp{
		List: respList,
	}, nil
}
