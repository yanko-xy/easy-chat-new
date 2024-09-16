package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/yanko-xy/easy-chat/apps/user/models"

	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"time"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendAlreadyExist = xerr.NewMsgErr("该用户已经是你的好友")
	ErrFriendRequestExist = xerr.NewMsgErr("该好友申请已存在")
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务：请求好友、通过或拒绝申请、好友列表
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {

	// 申请人是否与目标识好友关系
	friend, err := l.svcCtx.FriendModel.FindByUserIdAndFriendId(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend by userId and friendUid err %v, "+
			"req %v", err, in)
	}

	// 已经是好友
	if friend != nil {
		return nil, errors.WithStack(ErrFriendAlreadyExist)
	}

	// 是否已经申请
	friendReq, err := l.svcCtx.FriendRequestModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request by reqUid and userId err %v, "+
			"req %v", err, in)
	}
	if friendReq != nil {
		return nil, errors.WithStack(ErrFriendRequestExist)
	}

	// 创建申请记录
	_, err = l.svcCtx.FriendRequestModel.Insert(l.ctx, &socialmodels.FriendRequest{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: time.Now(),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandleResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friend request err %v, req %v", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
