package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/pkg/xcopy"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

// 定义两个结构体
type Source struct {
	Timestamp sql.NullTime
}

type Destination struct {
	Timestamp int64
}

type GroupPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInListLogic) GroupPutInList(in *social.GroupPutInListReq) (*social.GroupPutInListResp, error) {
	groupReqs, err := l.svcCtx.GroupRequestModel.ListNoHandle(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group put in by group id err %v, req %v",
			err, in.GroupId)
	}

	var respList []*social.GroupRequest
	xcopy.Copy(&respList, &groupReqs)

	return &social.GroupPutInListResp{
		List: respList,
	}, nil
}
