/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex
	opt           *serverOption
	routes        map[string]HandlerFunc
	authorization Authorization
	pattern       string
	userToConn    map[string]*Conn
	connToUser    map[*Conn]string
	addr          string
	upgrader      *websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		opt:           opt,
		routes:        make(map[string]HandlerFunc),
		userToConn:    make(map[string]*Conn),
		connToUser:    make(map[*Conn]string),
		addr:          addr,
		authorization: opt.Authorization,
		pattern:       opt.pattern,
		upgrader:      &websocket.Upgrader{},
		Logger:        logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	// 鉴权
	if !s.authorization.Auth(w, r) {
		s.Error("authorization failed")
		w.Write([]byte("authorization failed"))
		return
	}

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handleConn(conn)
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authorization.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登录过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// 根据连接对象执行任务处理
func (s *Server) handleConn(conn *Conn) {
	for {
		// 获取请求信息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}

		// 解析消息
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Errorf("websocket conn read message unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 依据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			// ping 回复
			s.Send(&Message{
				FrameType: FramePing,
			}, conn)
		case FrameData:
			// 根据请求的method分发路由并执行
			if handle, ok := s.routes[message.Method]; ok {
				handle(s, conn, &message)
			} else {
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在执行的方法 %v, 请检查", message.Method),
				}, conn)
			}
		}

	}
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}

	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []*Route) {
	for _, v := range rs {
		s.routes[v.Method] = v.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.pattern, s.ServerWs)
	s.Error(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止websocket服务")
}
