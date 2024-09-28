package group

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleReq) (resp *types.GroupPutInHandleResp, err error) {
	_, err = l.svcCtx.Social.GroupPutInHandle(l.ctx, &social.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleUid:    ctxdata.GetUId(l.ctx),
		HandleResult: req.HandleResult,
	})

	return
}
