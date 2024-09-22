package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/im/api/internal/types"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/imclient"
	"github.com/yanko-xy/easy-chat/pkg/ctxdata"
	"github.com/yanko-xy/easy-chat/pkg/wuid"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户获取聊天记录
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// 判断会话id是否是自己的
	ids := wuid.DeCombineId(req.ConversationId)
	uid := ctxdata.GetUId(l.ctx)
	if ids[0] != uid && ids[1] != uid {
		return nil, xerr.NewIllegalOperationErr()
	}

	data, err := l.svcCtx.Im.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId,
		StartSendTime:  req.StartSendTime,
		EndSendTime:    req.EndSendTime,
		Count:          req.Count,
		MsgId:          req.MsgId,
	})

	if err != nil {
		return nil, err
	}

	var res types.ChatLogResp
	copier.Copy(&res, data)

	return &res, err
}
