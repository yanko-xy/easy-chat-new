/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package svc

import "github.com/yanko-xy/easy-chat/apps/im/ws/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
