/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	*websocket.Conn
	Uid string
	s   *Server

	idleMu            sync.Mutex
	idle              time.Time
	maxConnectionIdle time.Duration

	messageMu       sync.Mutex
	readMessages    []*Message
	readMessagesSeq map[string]*Message
	message         chan *Message

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade http conn err %v", err)
		return nil
	}

	conn := &Conn{
		Conn: c,
		s:    s,

		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,

		readMessages:    make([]*Message, 0, 2),
		readMessagesSeq: make(map[string]*Message, 2),
		message:         make(chan *Message, 1),

		done: make(chan struct{}),
	}

	go conn.keepalive()
	return conn
}

// 长连接检测机制
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle

			fmt.Printf("idle %v, maxIdle %v", c.idle, c.maxConnectionIdle)
			if idle.IsZero() {
				// the connection is non-idle
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}

			val := c.maxConnectionIdle - time.Since(idle)
			fmt.Printf("val %v\n", val)
			c.idleMu.Unlock()
			if val < 0 {
				// the connection has been idle for a duration of keepalive.MaxConnectionIdle or more
				// Gracefully close the connection
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			fmt.Println("客户端结束连接")
			return
		}
	}
}

// 将消息记录到队列中
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 读队列中
	if m, ok := c.readMessagesSeq[msg.Id]; ok {
		// 已经有消息的记录，该消息已经有ack确认
		if len(c.readMessages) == 0 {
			// 队列中没有该消息
			return
		}

		// msg.AckSeq > m.AckSqe
		if m.Id != msg.Id || m.AckSeq >= msg.AckSeq {
			// 没有进行ack的确认，重复
			return
		}

		c.readMessagesSeq[msg.Id] = msg
		return
	}

	// 还没有进行ack的确认，避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessages = append(c.readMessages, msg)
	c.readMessagesSeq[msg.Id] = msg
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	// 开始忙碌
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	err := c.Conn.WriteMessage(messageType, data)
	// 当写操作完成后当前连接就会进入空闲状态，并记录空闲状态
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}
