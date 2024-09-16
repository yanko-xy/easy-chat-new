package logic

import (
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/config"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/internal/svc"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/conf"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
