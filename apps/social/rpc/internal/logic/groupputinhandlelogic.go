package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqBeforePass   = xerr.NewMsgErr("群申请请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsgErr("群申请请求已拒绝")
	ErrGroupReqBeforeCancel = xerr.NewMsgErr("群申请请求已取消")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	groupReq, err := l.svcCtx.GroupRequestModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group request by id err %v, req %v", err, in.GroupReqId)
	}

	handlerUser, err := l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.HandleUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by group id and user id err %v %v, req %v",
			err, in.HandleUid, in.GroupId)
	}
	// 操作人员是否有权限
	if constants.GroupRoleLevel(handlerUser.RoleLevel) == constants.AtLargeGroupRoleLevel {
		return nil, errors.WithStack(xerr.NewIllegalOperationErr())
	}

	switch constants.HandleResult(groupReq.HandleResult.Int64) {
	case constants.PassHandleResutl:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandleResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	case constants.CancelHandleResult:
		return nil, errors.WithStack(ErrGroupReqBeforeCancel)
	}

	// 更改群申请状态
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	groupReq.HandleTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	groupReq.HandleUserId = sql.NullString{
		String: in.HandleUid,
		Valid:  true,
	}

	err = l.svcCtx.GroupRequestModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestModel.Update(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update group request err %v, req %v", err, groupReq)
		}

		if constants.HandleResult(groupReq.HandleResult.Int64) != constants.PassHandleResutl {
			return nil
		}

		groupMember := &socialmodels.GroupMember{
			GroupId:   groupReq.GroupId,
			UserId:    groupReq.ReqId,
			RoleLevel: int64(constants.AtLargeGroupRoleLevel),
			JoinTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			JoinSource: sql.NullInt64{
				Int64: groupReq.JoinSource.Int64,
				Valid: true,
			},
			InviterUid: sql.NullString{
				String: groupReq.InviterUserId.String,
				Valid:  true,
			},
			OperatorUid: sql.NullString{
				String: in.HandleUid,
				Valid:  true,
			},
		}
		_, err = l.svcCtx.GroupMemberModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v, req %v", err, groupMember)
		}

		return nil
	})

	return &social.GroupPutInHandleResp{GroupId: groupReq.GroupId}, nil
}
