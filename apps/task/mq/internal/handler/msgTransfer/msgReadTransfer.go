/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/rogpeppe/go-internal/cache"
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/apps/task/mq/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/task/mq/mq"
	"github.com/yanko-xy/easy-chat/pkg/bitmap"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/zeromicro/go-queue/kq"
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	// 无
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseTransfer
	cache.Cache
	mu sync.Mutex

	groupMsgs map[string]*groupMsgRead
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseTransfer: NewBaseTransfer(svc),
		groupMsgs:    make(map[string]*groupMsgRead, 1),
		push:         make(chan *ws.Push, 1),
	}

	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	go m.transfer()

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Info("MsgReadTransfer", value)

	var data mq.MsgMarkRead

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 更新消息聊天记录中的已读状态
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	push := &ws.Push{
		ChatType:       data.ChatType,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ReadRecords:    readRecords,
		ContentType:    constants.ContentMakeRead,
	}

	switch data.ChatType {
	case constants.SingleChatType:
		// 直接推送
		m.push <- push
	case constants.GroupChatType:
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = ""

		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push %v", push.ConversationId)
			// 合并请求
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.groupMsgs[push.ConversationId] = NewGroupMsgRead(push, m.push)
		}
	}

	return nil
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return res, err
	}

	m.Infof("chatLogs %v", chatLogs)

	for _, chatlog := range chatLogs {
		switch data.ChatType {
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatlog.ReadRecords)
			readRecords.Set(data.SendId)
			chatlog.ReadRecords = readRecords.Export()
		case constants.SingleChatType:
			chatlog.ReadRecords = []byte{1}
		}

		res[chatlog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatlog.ReadRecords)

		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatlog.ID, chatlog.ReadRecords)
		if err != nil {
			m.Errorf("update make read err %v", err)
		}
	}

	return res, nil
}

func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("m transfer err %v , push %v", err, push)
			}
		}

		if push.ChatType == constants.SingleChatType {
			continue
		}

		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}

		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}
		m.mu.Unlock()
	}
}
