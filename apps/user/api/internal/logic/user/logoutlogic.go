package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/yanko-xy/easy-chat/apps/user/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退出登录
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err = l.svcCtx.Redis.HdelCtx(l.ctx, constants.REDIS_ONLINE_USER, uid)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "redis hdelctx %v err %v", constants.REDIS_ONLINE_USER, err)
	}
	return
}
