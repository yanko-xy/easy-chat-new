package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/im"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, immodels.ErrNotFound) {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversations by userId err %v, req %v", err, in.UserId)
	}

	var res im.GetConversationsResp
	copier.Copy(&res, &data)

	ids := make([]string, 0, len(data.ConversationList))
	for _, id := range data.ConversationList {
		ids = append(ids, id.ConversationId)
	}

	// 统计会话的消息情况
	list, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list by conversationIds err %v, req %v", err, ids)
	}

	for _, conversation := range list {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}

		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 有新消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 待读消息量
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 有新消息显示会话框
			res.ConversationList[conversation.ConversationId].IsShow = true
			// 最后一条消息
			if conversation.Msg != nil {
				res.ConversationList[conversation.ConversationId].Msg = &im.ChatLog{
					MsgType:    int32(conversation.Msg.MsgType),
					MsgContent: conversation.Msg.MsgContent,
				}
			}

		}
	}
	return &res, nil
}
