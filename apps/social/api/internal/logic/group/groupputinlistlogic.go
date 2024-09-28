package group

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/socialclient"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群列表
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListReq) (resp *types.GroupPutInListResp, err error) {
	list, err := l.svcCtx.Social.GroupPutInList(l.ctx, &socialclient.GroupPutInListReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, err
	}

	var respList []*types.GroupRequest
	copier.Copy(&respList, list.List)

	return &types.GroupPutInListResp{
		List: respList,
	}, nil
}
