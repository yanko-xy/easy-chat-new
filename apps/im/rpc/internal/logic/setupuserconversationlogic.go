package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/apps/im/immodels"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/im"
	"github.com/yanko-xy/easy-chat/apps/im/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/yanko-xy/easy-chat/pkg/wuid"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	var res im.SetUpUserConversationResp

	switch constants.ChatType(in.ChatType) {
	case constants.GroupChatType:
	// todo: 群聊
	case constants.SingleChatType:
		// 生成会话id
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 验证是否建立过会话
		_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 建立会话
			if errors.Is(err, immodels.ErrNotFound) {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation err %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation err %v, req %v", err, conversationId)
			}

		}

		err = l.setupUserConversation(conversationId, in.SendId, constants.SingleChatType, true)
		if err != nil {
			return &res, err
		}
		// 接收者是被动与目标用户建立连接，因此理论上是不需要在会话列表里展示，而是在用户发起聊天后展示
		err = l.setupUserConversation(conversationId, in.RecvId, constants.SingleChatType, false)
		if err != nil {
			return &res, err
		}
	}
	return &res, nil
}

func (l *SetUpUserConversationLogic) setupUserConversation(conversationId, userId string, chatType constants.ChatType, isShow bool) error {

	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		// 没有会话列表
		if errors.Is(err, immodels.ErrNotFound) {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
				CreateAt:         time.Now(),
			}

		} else {
			return err
		}

	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		// 存在
		return nil
	}

	// 不存在，需要建立
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       chatType,
		IsShow:         isShow,
	}

	// 存在更新，不存在新增
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "update conversations err %v", err)
	}
	return nil
}
