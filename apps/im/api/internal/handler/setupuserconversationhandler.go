package handler

import (
	"github.com/yanko-xy/easy-chat/pkg/validator"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"net/http"

	"github.com/yanko-xy/easy-chat/apps/im/api/internal/logic"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 建立会话
func setUpUserConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetUpUserConversationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.NewParameterErr("参数错误: "+err.Error()))
			return
		}

		// validator10 验证
		if errMsg, errCode := validator.Validate(req); errCode != 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.NewParameterErr(errMsg))
			return
		}

		l := logic.NewSetUpUserConversationLogic(r.Context(), svcCtx)
		resp, err := l.SetUpUserConversation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
