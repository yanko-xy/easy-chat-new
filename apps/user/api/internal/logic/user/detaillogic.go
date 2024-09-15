package user

import (
	"context"
	"rpc/user"

	"github.com/jinzhu/copier"
	"github.com/yanko-xy/easy-chat/apps/user/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/api/internal/types"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	userInfoResp, err := l.svcCtx.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	var user types.User
	copier.Copy(&user, userInfoResp.User)

	return &types.UserInfoResp{
		Info: user,
	}, nil
}
