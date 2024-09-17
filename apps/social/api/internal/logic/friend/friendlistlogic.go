package friend

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/userclient"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	friends, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	if len(friends.List) == 0 {
		return nil, nil
	}

	// 根据好友id数组
	uids := make([]string, 0, len(friends.List))
	for _, v := range friends.List {
		uids = append(uids, v.FriendUid)
	}

	// 根据好友id数组查询好友信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return nil, err
	}

	userRecords := make(map[string]*userclient.UserEntity, len(users.Users))
	for _, v := range users.Users {
		userRecords[v.Id] = v
	}

	respList := make([]*types.Friend, 0, len(friends.List))
	for _, v := range friends.List {
		friend := &types.Friend{
			Id:        int64(v.Id),
			FriendUid: v.FriendUid,
			Remark:    v.Remark,
		}

		if u, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}

		respList = append(respList, friend)
	}

	return &types.FriendListResp{
		List: respList,
	}, nil
}
