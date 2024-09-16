package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &social.GroupPutInResp{}, nil
}
