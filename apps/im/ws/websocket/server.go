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
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	default:

	}

	return "NoAck"
}

type Server struct {
	sync.RWMutex
	opt *serverOption
	*threading.TaskRunner
	routes        map[string]HandlerFunc
	authorization Authorization
	pattern       string
	userToConn    map[string]*Conn
	connToUser    map[*Conn]string
	addr          string
	upgrader      *websocket.Upgrader

	listenOn string
	discover Discover
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	s := &Server{
		opt:           opt,
		TaskRunner:    threading.NewTaskRunner(opt.concurrency),
		routes:        make(map[string]HandlerFunc),
		userToConn:    make(map[string]*Conn),
		connToUser:    make(map[*Conn]string),
		addr:          addr,
		authorization: opt.Authorization,
		pattern:       opt.pattern,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		discover: opt.discover,
		listenOn: FigureOutListenOn(addr),
		Logger:   logx.WithContext(context.Background()),
	}

	// 存在服务发现，采用分布式im通信的时候；默认不做任何处理
	s.discover.Register(fmt.Sprintf("%s", s.listenOn))

	return s
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

// 根据连接对象执行任务处理
func (s *Server) handleConn(conn *Conn) {

	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 如果存在服务发现则进行注册；默认不做任何处理
	s.discover.BindUser(conn.Uid)

	// 处理任务
	go s.hadnlerWrite(conn)

	if s.isAck(nil) {
		go s.readAck(conn)
	}

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

			//s.Close(conn)
			//return
			continue
		}

		// todo: 给客户端回复一个ack

		// 依据请求消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}

	}
}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck && message.FrameType != FrameTranspond
}

// 读取消息的ack
func (s *Server) readAck(conn *Conn) {

	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			return nil
		}

		s.Errorf("message ack send err %v message %v", err, msg)
		conn.readMessages[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessages[0].errCount) * time.Microsecond
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}

		time.Sleep(tempDelay)
		return err
	}

	for {
		select {
		case <-conn.done:
			// 关闭连接
			s.Infof("close message ack uid %v", conn.Uid)
			return
		default:

		}

		conn.messageMu.Lock()

		// 从列表中读取新的消息
		if len(conn.readMessages) == 0 {
			conn.messageMu.Unlock()
			// 增加睡眠
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 读取第一条
		message := conn.readMessages[0]

		// 判断ack的方式
		switch s.opt.ack {
		case OnlyAck:
			// 直接给客户端回复
			if err := send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn); err != nil {
				continue
			}
			// 进行业务处理
			// 把消息从队列中移除
			conn.readMessages = conn.readMessages[1:]
			conn.messageMu.Unlock()
			conn.message <- message
		case RigorAck:
			// 先回
			if message.AckSeq == 0 {
				// 还未确认
				conn.readMessages[0].AckSeq++
				conn.readMessages[0].ackTime = time.Now()
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn); err != nil {
					continue
				}
				s.Infof("message ack RigorAck send mid %v, seq %v, time %v", message.Id, message.AckSeq, message.ackTime)
				conn.messageMu.Unlock()
				time.Sleep(300 * time.Millisecond)
				continue
			}

			// 再验证

			// 1. 客户端返回结果，再一次确认
			// 得到客户端的序号
			msgSeq := conn.readMessagesSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 确认
				conn.readMessages = conn.readMessages[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 2.1 超过确认时间
				s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
				delete(conn.readMessagesSeq, message.Id)
				conn.readMessages = conn.readMessages[1:]
				conn.messageMu.Unlock()
				continue
			}

			// 2.2 未超时, 重新发送
			conn.messageMu.Unlock()
			if err := send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn); err != nil {
				continue
			}
			// 睡眠一定时间
			time.Sleep(300 * time.Millisecond)

		default:

		}
	}
}

// 任务的处理
func (s *Server) hadnlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 关闭连接
			return
		case message := <-conn.message:
			// 根据请求的method分发路由并执行
			switch message.FrameType {
			case FramePing:
				// ping 回复
				s.Send(&Message{
					FrameType: FramePing,
				}, conn)
			case FrameData:
				// 根据请求的method分发路由并执行
				if handle, ok := s.routes[message.Method]; ok {
					handle(s, conn, message)
				} else {
					s.Send(&Message{
						FrameType: FrameData,
						Data:      fmt.Sprintf("不存在执行的方法 %v, 请检查", message.Method),
					}, conn)
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessagesSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
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

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) ([]*Conn, []string) {
	if len(uids) == 0 {
		return nil, nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	noExistUids := make([]string, 0)
	for _, uid := range uids {
		if conn, ok := s.userToConn[uid]; ok {
			res = append(res, conn)
		} else {
			noExistUids = append(noExistUids, uid)
		}
	}

	return res, noExistUids
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

	s.discover.RelieveUser(uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	conns, noExistUids := s.GetConns(sendIds...)

	// 发送当前服务存在的用户
	err := s.Send(msg, conns...)
	if err != nil {
		return err
	}

	// 不存在的，转发由其他服务处理
	return s.discover.Transpond(msg, noExistUids...)
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
