package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	friendList, err := l.svcCtx.FriendModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "friend list by userId err %v, req %v", err, in)
	}

	var respList []*social.Friend
	copier.Copy(&respList, &friendList)
	return &social.FriendListResp{
		List: respList,
	}, nil
}
