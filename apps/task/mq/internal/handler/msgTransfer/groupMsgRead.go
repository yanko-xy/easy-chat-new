/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package msgTransfer

import (
	"github.com/yanko-xy/easy-chat/apps/im/ws/ws"
	"github.com/yanko-xy/easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu             sync.Mutex
	conversationId string
	push           *ws.Push
	pushCh         chan *ws.Push
	count          int

	// 上次推送时间
	pushTime time.Time
	done     chan struct{}
}

func NewGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	go m.transfer()

	return m
}

func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

func (m *groupMsgRead) transfer() {
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			m.mu.Lock()

			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := m.push
			logx.Infof("timer.C %v val %v", time.Now(), val)
			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}

				// 未达标
				m.mu.Unlock()
				continue
			}

			pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()

			// 推送
			logx.Infof("超过 合并的条件推送 %v ", push)
			m.pushCh <- push
		case <-m.done:
			return
		default:
			m.mu.Lock()
			logx.Infof("groupMsgRead count %v, push %v", m.count, m.push)

			if m.count >= GroupMsgReadRecordDelayCount {
				logx.Infof("groupMsgRead transfer ConversationId %v , count %v", m.push.ConversationId, m.count)
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()

				logx.Infof("达到推送量推送 %v", push)
				m.pushCh <- push
				continue
			}

			if m.isIdle() {
				m.mu.Unlock()
				// 使得MsgReadTransfer清理
				m.pushCh <- &ws.Push{
					ConversationId: m.conversationId,
					ChatType:       constants.GroupChatType,
				}
				continue
			}

			m.mu.Unlock()
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)

	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}

	return false
}

func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}
