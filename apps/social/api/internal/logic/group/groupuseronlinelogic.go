package group

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/socialclient"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, err
	}

	// 或取用户信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
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

	return &types.GroupUserOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
