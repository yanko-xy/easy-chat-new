package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"time"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInLogic) GroupPutIn(in *social.GroupPutInReq) (*social.GroupPutInResp, error) {

	//  1. 普通用户申请 ： 如果群无验证直接进入
	//  2. 群成员邀请： 如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		inviteGroupMember *socialmodels.GroupMember
		userGroupMember   *socialmodels.GroupMember
		groupInfo         *socialmodels.Group
		err               error
	)

	userGroupMember, err = l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.ReqUid, in.GroupId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by group id and user id err %v %v, req %v",
			err, in.ReqUid, in.GroupId)
	}
	// 已经在群里
	if userGroupMember != nil {
		return &social.GroupPutInResp{}, nil
	}

	groupReq, err := l.svcCtx.GroupRequestModel.FindByGroupIdAndReqUid(l.ctx, in.GroupId, in.ReqUid)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group request by group id and req uid err %v, req %v, %v",
			err, in.GroupId, in.ReqUid)
	}
	// 已经申请过
	if groupReq != nil {
		return &social.GroupPutInResp{}, nil
	}

	groupReq = &socialmodels.GroupRequest{
		ReqId:   in.ReqUid,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandleResult),
			Valid: true,
		},
	}

	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}

	groupInfo, err = l.svcCtx.GroupModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by id err %v, req %v", err, in.GroupId)
	}

	// 验证群是否需要验证
	if !groupInfo.IsVerify {
		// 不需要验证
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandleResutl),
			Valid: true,
		}
		groupReq.HandleTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		return l.createGroupReq(groupReq, true)
	}

	// 验证进群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 申请进入
		return l.createGroupReq(groupReq, false)
	}

	inviteGroupMember, err = l.svcCtx.GroupMemberModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by group id and user id err %v, req %v %v", err, in.InviterUid, in.GroupId)
	}
	// 判断邀请人身份
	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel ||
		constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 管理员/群创建者邀请
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandleResutl),
			Valid: true,
		}
		groupReq.HandleTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	// 普通成员邀请
	return l.createGroupReq(groupReq, false)
}

func (l *GroupPutInLogic) createGroupReq(groupReq *socialmodels.GroupRequest, isPass bool) (*social.GroupPutInResp, error) {
	_, err := l.svcCtx.GroupRequestModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group request err %v, req %v", err, groupReq)
	}

	if isPass {
		return &social.GroupPutInResp{GroupId: groupReq.GroupId}, nil
	}

	return &social.GroupPutInResp{}, nil
}

func (l *GroupPutInLogic) createGroupMember(in *social.GroupPutInReq) error {
	groupMember := &socialmodels.GroupMember{
		GroupId:   in.GroupId,
		UserId:    in.ReqUid,
		RoleLevel: int64(constants.AtLargeGroupRoleLevel),
		JoinTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		OperatorUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
	}
	_, err := l.svcCtx.GroupMemberModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v, req %v", err, groupMember)
	}

	return nil
}
