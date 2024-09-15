package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/user"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = xerr.New(xerr.SERVER_COMMON_ERROR, "该用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {

	userEntity, err := l.svcCtx.UserModels.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrUserNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id err %v, req %v", err, in.Id)
	}

	var resp user.UserEntity
	copier.Copy(&resp, userEntity)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
