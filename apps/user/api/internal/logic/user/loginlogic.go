package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/yanko-xy/easy-chat/apps/user/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	loginResp, err := l.svcCtx.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 设置用户在线
	err = l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "redis hsetctx %v err %v", constants.REDIS_ONLINE_USER, err)
	}

	return &res, nil
}
