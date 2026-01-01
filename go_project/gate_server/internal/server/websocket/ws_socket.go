package websocket

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mxxmstar/learning/gate_server/internal/conn"
	auth_user "github.com/mxxmstar/learning/gate_server/internal/user_auth"
	"github.com/mxxmstar/learning/pkg/logger"
	"go.uber.org/zap"
)

// MessageHandler 定义消息处理器接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, conn conn.Connection, envelope *Envelope) error
}

// client 与 server 通信的消息格式
type Envelope struct {
	Type      string                 `json:"type"`
	Token     string                 `json:"token,omitempty"`
	SessionId string                 `json:"session_id,omitempty"`
	DeviceId  string                 `json:"device_id,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"` // 消息体 方便扩展
}

// wsConnection 实现 conn.Connection 接口
type wsConnection struct {
	connId    string
	userId    uint64
	deviceId  string
	ws        *websocket.Conn
	SendChan  chan []byte            // 发送消息的 channel
	closed    atomic.Bool            // 是否关闭
	closeChan chan struct{}          // 关闭 channel
	mgr       conn.ConnectionManager // 连接管理器
}

func (c *wsConnection) Id() string {
	return c.connId
}
func (c *wsConnection) UserId() uint64 {
	return c.userId
}
func (c *wsConnection) Send(msg []byte) error {
	if c.closed.Load() {
		return conn.ErrConnectionClosed
	}

	select {
	case c.SendChan <- msg:
		return nil
	case <-c.closeChan:
		return conn.ErrConnectionClosed
	default:
		return conn.ErrConnectionClosed
	}
}

func (c *wsConnection) Close(reason string) error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil
	}
	close(c.closeChan)
	logger.FormatLog(context.Background(), "info", "wsConnection Close",
		zap.String("connId", c.Id()),
		zap.Uint64("userId", c.UserId()),
		zap.String("reason", reason))

	// 发送控制帧关闭连接
	_ = c.ws.WriteControl(
		websocket.CloseMessage, // 关闭消息
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason), // 格式化关闭消息体内容(正常关闭1000)
		time.Now().Add(time.Second),                                        // 添加超时时间
	)
	_ = c.ws.Close()
	c.mgr.UnRegister(c)
	return nil
}

const (
	authTimeout     = 5 * time.Second  // 验证 token/session 超时时间 (ValidateTokenOrSession)
	sessionTTl      = 300              // session 过期时间
	pongwait        = 60 * time.Second // 读超时
	pingPeriod      = 25 * time.Second // ping 间隔
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// NotifyOldFunc 通知旧连接关闭
type NotifyOldFunc func(ctx context.Context, oldConnId string)

type WebsocketServer struct {
	gateId    string
	mgr       conn.ConnectionManager
	auth      auth_user.AuthService     // 验证服务
	handlers  map[string]MessageHandler // 消息处理器映射
	store     interface{}               // TODO:会话管理，连接状态存储，踢掉旧连接，用户在线状态管理，分布式存储
	notifyOld NotifyOldFunc
	upgrader  websocket.Upgrader
}

func NewWebsocketServer(
	gateId string,
	mgr conn.ConnectionManager,
	auth auth_user.AuthService,
	store interface{},
	notifyOld NotifyOldFunc,
) *WebsocketServer {
	s := &WebsocketServer{
		gateId:    gateId,
		mgr:       mgr,
		auth:      auth,
		handlers:  make(map[string]MessageHandler),
		store:     store,
		notifyOld: notifyOld,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: 生产环境限制 Origin
			},
		},
	}
	return s
}

// 注册消息处理器,由消息处理器处理各种业务消息
func (s *WebsocketServer) RegisterHandler(msgType string, handler MessageHandler) {
	s.handlers[msgType] = handler
}

// 处理新连接和初始认证
func (s *WebsocketServer) handleNewConnectioon(w http.ResponseWriter, r *http.Request) {
	// 升级为 websocket
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] upgrade failed: %v", err))
		return
	}
	defer ws.Close()

	// 读取并处理认证消息
	_ = ws.SetReadDeadline(time.Now().Add(authTimeout)) // 设置读超时
	_, msg, err := ws.ReadMessage()
	if err != nil {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] read auth message failed: %v", err))
		return
	}

	var envelope Envelope
	if err := json.Unmarshal(msg, &envelope); err != nil || envelope.Type != "auth" {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] invalid auth message: %v", err))
		_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_nack","reason":"invalid format"}`))
		return
	}

	// 执行认证
	ctx, cancel := context.WithTimeout(r.Context(), authTimeout)
	authResult, err := s.auth.ValidateTokenOrSession(ctx, envelope.Token, envelope.SessionId, envelope.DeviceId)
	cancel()
	if err != nil || !authResult.Valid {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] auth failed: %v, error: %s", err, authResult.Error))
		_ = ws.WriteMessage(websocket.TextMessage, []byte(`"type":"auth_nack","reason":"auth failed"}`))
		return
	}

	// 创建连接对象并注册
	// connId := logger.NewTraceId()
	connId := fmt.Sprintf("%s#%s", s.gateId, s.randUUID())
	wsConn := &wsConnection{
		connId:    connId,
		userId:    authResult.UserId,
		deviceId:  authResult.DeviceId,
		ws:        ws,
		SendChan:  make(chan []byte, 256),
		closeChan: make(chan struct{}),
		mgr:       s.mgr,
	}

	if err := s.mgr.Register(wsConn); err != nil {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] register failed: %v", err))
		_ = ws.WriteMessage(websocket.TextMessage, []byte(`"type":"auth_nack","reason":"register ws connection failed"}`))
		return
	}

	// 发送认证成功响应 TODO
	ack := map[string]interface{}{
		"type":        "auth_ack",
		"conn_id":     connId,
		"session_ttl": sessionTTl,
		"server_time": time.Now().Unix(),
		// "body": map[string]interface{}{
		// 	"conn_id": connId,
		// },
	}
	ackBytes, _ := json.Marshal(ack)
	if err := ws.WriteMessage(websocket.TextMessage, ackBytes); err != nil {
		logger.FormatLog(r.Context(), "error", fmt.Sprintf("[ws] write auth ack failed: %v", err))
		wsConn.Close("auth ack failed")
		return
	}

	// 启动消息处理协程
	go s.readPump(wsConn)
	go s.writePump(wsConn)
}

func (s *WebsocketServer) readPump(wsConn *wsConnection) {
	defer wsConn.Close("read pump exit")
	ws := wsConn.ws
	ws.SetReadLimit(16 * 1024)
	_ = ws.SetReadDeadline(time.Now().Add(pongwait))
	// 设置心跳处理器
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongwait))
		return nil
	})

	for {
		select {
		case <-wsConn.closeChan:
			// 连接被关闭
			return
		default:
		}

		// 读取消息
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 非正常错误，记录日志
				logger.FormatLog(context.Background(), "error", fmt.Sprintf("[conn %s] read error: %v", wsConn.connId, err))
			}
			// 遇到错误，退出主循环关闭连接
			return
		}

		// 解析消息并路由到相应的处理器
		var envelope Envelope
		if err := json.Unmarshal(msg, &envelope); err != nil {
			logger.FormatLog(context.Background(), "error", fmt.Sprintf("[conn %s] unmarshal message error: %v", wsConn.connId, err))
			continue
		}

		// 特殊处理 ping 消息
		if envelope.Type == "ping" {
			_ = wsConn.Send([]byte(`{"type":"pong"}`))
			continue
		}

		// 查找并执行对应的消息处理器
		if handler, exists := s.handlers[envelope.Type]; exists {
			ctx := context.WithValue(context.Background(), "conn", wsConn)
			if err := handler.HandleMessage(ctx, wsConn, &envelope); err != nil {
				logger.FormatLog(context.Background(), "error", fmt.Sprintf("[conn %s] handle message error: %v", wsConn.connId, err))
			}
		} else {
			// 未知消息类型，可以选择忽略或记录日志
			logger.FormatLog(context.Background(), "warn", fmt.Sprintf("[conn %s] unknown message type: %s", wsConn.connId, envelope.Type))
		}
	}
}

func (s *WebsocketServer) writePump(wsConn *wsConnection) {
	// 创建一个定时器，定时发送ping消息
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wsConn.Close("write pump exit")
	}()

	ws := wsConn.ws
	for {
		select {
		case <-wsConn.closeChan:
			// 关闭连接
			return
		case msg, ok := <-wsConn.SendChan:
			// 发送消息请求
			if !ok {
				_ = ws.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.FormatLog(context.Background(), "error", fmt.Sprintf("[conn %s] write error: %v", wsConn.connId, err))
				return
			}
		case <-ticker.C:
			// ping 心跳消息
			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.FormatLog(context.Background(), "error", fmt.Sprintf("[conn %s] ping error: %v", wsConn.connId, err))
				return
			}
		}
	}
}

// 工具函数
func (s *WebsocketServer) randUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
