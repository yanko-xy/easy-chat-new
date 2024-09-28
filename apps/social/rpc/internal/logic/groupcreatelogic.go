package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/apps/social/socialmodels"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/wuid"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群业务：创建群，修改群，群公告，申请群，用户群列表，群成员，申请群，群退出..
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {

	group := &socialmodels.Group{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		IsVerify:   in.IsVerify,
	}

	err := l.svcCtx.GroupModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupModel.Insert(l.ctx, session, group)

		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v, req %v", err, group)
		}

		groupMember := &socialmodels.GroupMember{
			GroupId:   group.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int64(constants.CreatorGroupRoleLevel),
			JoinTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}
		_, err = l.svcCtx.GroupMemberModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v, req %v", err, groupMember)
		}
		return nil
	})

	return &social.GroupCreateResp{
		GroupId: group.Id,
	}, err
}
