package logic

import (
	"context"

	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/user/models"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {

	var (
		userEntitys []*models.User
		err         error
	)

	if in.Phone != "" {
		userEntity, err := l.svcCtx.UserModels.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		userEntitys, err = l.svcCtx.UserModels.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntitys, err = l.svcCtx.UserModels.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user err %v, req %v", err, in)
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntitys)

	return &user.FindUserResp{
		Users: resp,
	}, nil
}
