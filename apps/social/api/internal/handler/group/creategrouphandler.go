package group

import (
	"github.com/yanko-xy/easy-chat/pkg/validator"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"net/http"

	"github.com/yanko-xy/easy-chat/apps/social/api/internal/logic/group"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创群
func CreateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.NewParameterErr("参数错误: "+err.Error()))
			return
		}

		// validator10 验证
		if errMsg, errCode := validator.Validate(req); errCode != 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.NewParameterErr(errMsg))
			return
		}

		l := group.NewCreateGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
