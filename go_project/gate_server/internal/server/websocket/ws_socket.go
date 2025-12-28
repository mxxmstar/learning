package websocket

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mxxmstar/learning/gate_server/internal/conn"
	"github.com/mxxmstar/learning/pkg/logger"
	"go.uber.org/zap"
)

// client 与 server 通信的消息格式
type Envelope struct {
	Type      string                 `json:"type"`
	Token     string                 `json:"token,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"` // 消息体 方便扩展
}

// wsConnection 实现 conn.Connection 接口
type wsConnection struct {
	connID    string
	userID    uint64
	ws        *websocket.Conn
	SendChan  chan []byte            // 发送消息的 channel
	closed    atomic.Bool            // 是否关闭
	closeChan chan struct{}          // 关闭 channel
	mgr       conn.ConnectionManager // 连接管理器
}

func (c *wsConnection) ID() string {
	return c.connID
}
func (c *wsConnection) UserID() uint64 {
	return c.userID
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
		zap.String("connID", c.ID()),
		zap.Uint64("userID", c.UserID()),
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

// const (
// 	authTimeout     = 5 * time.Second  // 验证 token 超时时间
// 	sessionTTl      = 300              // session 过期时间
// 	pongwait        = 60 * time.Second // 读超时
// 	pingPeriod      = 25 * time.Second // ping 间隔
// 	readBufferSize  = 1024
// 	writeBufferSize = 1024
// )

// NotifyOldFunc 通知旧连接关闭
type NotifyOldFunc func(ctx context.Context, oldConnID string)

// type websocketServer struct {
// 	gateID   string
// 	mgr      conn.ConnectionManager
// 	verifier conn.Verifier
//     store     *storage.Store
// 	notifyOld NotifyOldFunc
// 	upgrader  websocket.Upgrader
// 	srvMux    *http.ServeMux
// }
// func NewWebsocketServer(
// 	gateID string,
// 	mgr conn.ConnectionManager,
// 	verifier auth.Verifier,
// 	store *storage.Store,
// 	notifyOld NotifyOldFunc,
// ) *WebsocketServer {
// 	s := &WebsocketServer{
// 		gateID:    gateID,
// 		mgr:       mgr,
// 		verifier:  verifier,
// 		store:     store,
// 		notifyOld: notifyOld,
// 		upgrader: websocket.Upgrader{
// 			ReadBufferSize:  readBufferSize,
// 			WriteBufferSize: writeBufferSize,
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true // TODO: 生产环境限制 Origin
// 			},
// 		},
// 	}
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/ws", s.handleWS)
// 	s.srvMux = mux
// 	return s
// }

// func (s *WebsocketServer) Router() http.Handler {
// 	return s.srvMux
// }

// func (s *WebsocketServer) handleWS(w http.ResponseWriter, r *http.Request) {
// 	ws, err := s.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("[ws] upgrade failed: %v", err)
// 		return
// 	}
// 	defer ws.Close()

// 	_ = ws.SetReadDeadline(time.Now().Add(authTimeout))
// 	_, msg, err := ws.ReadMessage()
// 	if err != nil {
// 		log.Printf("[ws] read auth message failed: %v", err)
// 		return
// 	}

// 	var env Envelope
// 	if err := json.Unmarshal(msg, &env); err != nil || env.Type != "auth" {
// 		log.Printf("[ws] invalid auth frame")
// 		_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_nack","reason":"invalid format"}`))
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
// 	ar, err := s.verifier.ValidateToken(ctx, env.Token, env.DeviceID)
// 	cancel()
// 	if err != nil {
// 		log.Printf("[ws] auth failed: %v", err)
// 		_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_nack","reason":"auth failed"}`))
// 		return
// 	}

// 	connID := fmt.Sprintf("%s#%s", s.gateID, randUUID())
// 	oldConn, err := s.store.SwapSessionAtomic(r.Context(), ar.UserID, env.DeviceID, connID, sessionTTL)
// 	if err != nil {
// 		log.Printf("[ws] session swap failed: %v", err)
// 		_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_nack","reason":"session error"}`))
// 		return
// 	}

// 	if oldConn != "" && oldConn != connID {
// 		log.Printf("[ws] old connection %s exists for user %d device %s", oldConn, ar.UserID, env.DeviceID)
// 		if s.notifyOld != nil {
// 			go func() {
// 				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 				defer cancel()
// 				s.notifyOld(ctx, oldConn)
// 			}()
// 		}
// 	}

// 	wc := &wsConnection{
// 		connID:  connID,
// 		userID:  ar.UserID,
// 		ws:      ws,
// 		sendCh:  make(chan []byte, 128),
// 		closeCh: make(chan struct{}),
// 		mgr:     s.mgr,
// 	}

// 	if err := s.mgr.Register(wc); err != nil {
// 		log.Printf("[ws] register failed: %v", err)
// 		_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_nack","reason":"internal"}`))
// 		return
// 	}

// 	ack := map[string]interface{}{
// 		"type":        "auth_ok",
// 		"conn_id":     connID,
// 		"session_ttl": sessionTTL,
// 		"server_time": time.Now().Unix(),
// 	}
// 	ackb, _ := json.Marshal(ack)
// 	if err := ws.WriteMessage(websocket.TextMessage, ackb); err != nil {
// 		log.Printf("[ws] send ack failed: %v", err)
// 		wc.Close("write ack failed")
// 		return
// 	}

// 	go s.readPump(wc)
// 	go s.writePump(wc)
// }

// // --- Pumps ---

// func (s *WebsocketServer) readPump(wc *wsConnection) {
// 	defer wc.Close("read pump exit")
// 	ws := wc.ws
// 	ws.SetReadLimit(16 * 1024)
// 	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
// 	ws.SetPongHandler(func(string) error {
// 		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
// 		return nil
// 	})

// 	for {
// 		select {
// 		case <-wc.closeCh:
// 			return
// 		default:
// 		}

// 		_, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("[conn %s] read error: %v", wc.connID, err)
// 			}
// 			return
// 		}

// 		var m map[string]interface{}
// 		if err := json.Unmarshal(msg, &m); err == nil {
// 			if t, ok := m["type"].(string); ok {
// 				switch t {
// 				case "ping":
// 					_ = wc.Send([]byte(`{"type":"pong"}`))
// 					continue
// 				}
// 			}
// 		}
// 		// 默认：忽略未知消息（或记录）
// 	}
// }

// func (s *WebsocketServer) writePump(wc *wsConnection) {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		wc.Close("write pump exit")
// 	}()

// 	ws := wc.ws
// 	for {
// 		select {
// 		case <-wc.closeCh:
// 			return
// 		case msg, ok := <-wc.sendCh:
// 			if !ok {
// 				_ = ws.WriteMessage(websocket.CloseMessage, nil)
// 				return
// 			}
// 			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
// 			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
// 				log.Printf("[conn %s] write error: %v", wc.connID, err)
// 				return
// 			}
// 		case <-ticker.C:
// 			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
// 			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				log.Printf("[conn %s] ping error: %v", wc.connID, err)
// 				return
// 			}
// 		}
// 	}
// }

// // 工具函数
// func randUUID() string {
// 	b := make([]byte, 16)
// 	_, _ = rand.Read(b)
// 	return fmt.Sprintf("%x", b)
// }
