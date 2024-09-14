package logic

import (
	"context"
	"errors"
	"time"

	"rpc/internal/svc"
	"rpc/user"

	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/yanko-xy/easy-chat/pkg/encrypt"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = errors.New("手机号未注册")
	ErrUserPwdError     = errors.New("密码错误")
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
			return nil, ErrPhoneNotRegister
		}
		return nil, err
	}

	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, ErrUserPwdError
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
