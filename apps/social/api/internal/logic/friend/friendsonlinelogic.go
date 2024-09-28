package friend

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/socialclient"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	firendList, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	if len(firendList.List) == 0 {
		return &types.FriendsOnlineResp{}, nil
	}

	// 获取用户信息
	uids := make([]string, 0, len(firendList.List))
	for _, friend := range firendList.List {
		uids = append(uids, friend.FriendUid)
	}

	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "redis hgetall %v err %v", constants.REDIS_ONLINE_USER, err)
	}

	resOnlineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		if _, ok := onlines[s]; ok {
			resOnlineList[s] = true
		} else {
			resOnlineList[s] = false
		}
	}

	return &types.FriendsOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
