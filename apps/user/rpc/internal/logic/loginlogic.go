package logic

import (
	"context"
	"time"

	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"

	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/yanko-xy/easy-chat/pkg/encrypt"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号未注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 登录

	// 验证用户是否注册，根据手机号验证
	userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "gen token err %v, req %v", err, userEntity)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		Id:     userEntity.Id,
	}, nil
}
