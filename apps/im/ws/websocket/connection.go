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
	idlemu sync.Mutex
	*websocket.Conn
	s                 *Server
	idle              time.Time
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade http conn err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}),
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
			c.idlemu.Lock()
			idle := c.idle

			fmt.Printf("idle %v, maxIdle %v", c.idle, c.maxConnectionIdle)
			if idle.IsZero() {
				// the connection is non-idle
				c.idlemu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}

			val := c.maxConnectionIdle - time.Since(idle)
			fmt.Printf("val %v\n", val)
			c.idlemu.Unlock()
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

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	// 开始忙碌
	messageType, p, err = c.Conn.ReadMessage()
	c.idlemu.Lock()
	defer c.idlemu.Unlock()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	err := c.Conn.WriteMessage(messageType, data)
	// 当写操作完成后当前连接就会进入空闲状态，并记录空闲状态
	c.idlemu.Lock()
	defer c.idlemu.Unlock()
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
