package group

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInReq) (resp *types.GroupPutInResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	_, err = l.svcCtx.Social.GroupPutIn(l.ctx, &social.GroupPutInReq{
		GroupId:    req.GroupId,
		ReqUid:     uid,
		ReqMsg:     req.ReqMsg,
		JoinSource: int32(req.JoinSource),
		InviterUid: req.InviterUid,
	})
	return
}
