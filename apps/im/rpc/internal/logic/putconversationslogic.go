package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/im"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversations by userId err %v, req %v", err, in.UserId)
	}

	// 没有会话列表
	if conversations.ConversationList == nil {
		conversations.ConversationList = make(map[string]*immodels.Conversation)
	}

	for k, v := range in.ConversationList {
		var oldTotal int
		if conversations.ConversationList[k] != nil {
			oldTotal = conversations.ConversationList[k].Total
		}

		conversations.ConversationList[k] = &immodels.Conversation{
			ConversationId: v.ConversationId,
			ChatType:       constants.ChatType(v.ChatType),
			IsShow:         v.IsShow,
			Total:          int(v.Read) + oldTotal,
			Seq:            v.Seq,
		}
	}

	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update conversations err %v, req %v", err, conversations)
	}
	return &im.PutConversationsResp{}, nil
}
