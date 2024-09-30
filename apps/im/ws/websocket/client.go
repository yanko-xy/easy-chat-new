/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error

	Send(v any) error
	SendUid(v any, uids ...string) error
	Read(v any) error
}

type client struct {
	*websocket.Conn
	host string
	opt  *dailOption
	Discover
}

func NewClient(host string, opts ...DailOptins) *client {
	opt := newDailOption(opts...)

	c := &client{
		Conn:     nil,
		host:     host,
		opt:      opt,
		Discover: opt.Discover,
	}

	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return c
}

func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) SendUid(v any, uids ...string) error {
	if c.Discover != nil {
		return c.Discover.Transpond(v, uids...)
	}
	return c.Send(v)
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}

	// 增加一个重连发送
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
