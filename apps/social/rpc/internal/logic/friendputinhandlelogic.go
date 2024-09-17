package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsgErr("好友申请已经通过")
	ErrFirendReqBeforeRefuse = xerr.NewMsgErr("好友申请已经拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend req by id friendreqId err %v, "+
			"req %v", err, in.FriendReqId)
	}

	// 验证是否已经处理
	switch constants.HandleResult(friendReq.HandleResult.Int64) {
	case constants.PassHandleResutl:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandleResult:
		return nil, errors.WithStack(ErrFirendReqBeforeRefuse)
	}

	// 判断是否是本人操作
	if friendReq.UserId != in.UserId {
		return nil, errors.WithStack(xerr.NewIllegalOperationErr())
	}

	friendReq.HandleResult.Int64 = int64(in.HandleResult)
	friendReq.HandledAt = sql.NullTime{
		Time: time.Now(), Valid: true,
	}
	friendReq.HandleMsg = sql.NullString{
		String: in.HandleMsg,
		Valid:  true,
	}

	// 修改申请结果
	err = l.svcCtx.FriendRequestModel.Trans(l.ctx, func(ctx context.Context,
		session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestModel.UpdateWithTrans(l.ctx, session,
			friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend req err %v, req %v", err, friendReq)
		}

		if constants.HandleResult(in.HandleResult) != constants.PassHandleResutl {
			return nil
		}

		friends := []*socialmodels.Friend{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}

		_, err := l.svcCtx.FriendModel.BatchInsertWithTrans(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "friend batch insert err %v, req %v", err, friends)
		}

		return nil
	})

	return &social.FriendPutInHandleResp{}, err
}
