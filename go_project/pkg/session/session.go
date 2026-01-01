package session

import (
	"errors"
	"sync"
	"time"
)

// SessionType 会话类型
type SessionType string

const (
	UserStateSession     SessionType = "user_state"
	ClientGatewaySession SessionType = "client_gateway"
	LoginTokenSession    SessionType = "login_token"
	WebSession           SessionType = "web"
	IMAuthSession        SessionType = "im_auth"
)

var ErrJWTManagerNotSet = errors.New("jwt manager not set")

// SessionId 会话Id,唯一标识符
type SessionId string

type BaseSession struct {
	Id        SessionId
	Type      SessionType
	CreatedAt time.Time
	ExpiredAt time.Time
	Data      map[string]interface{}
	mutex     sync.RWMutex
}

func NewBaseSession(id SessionId, sessionType SessionType, ttl time.Duration) *BaseSession {
	now := time.Now()
	return &BaseSession{
		Id:        id,
		Type:      sessionType,
		CreatedAt: now,
		ExpiredAt: now.Add(ttl),
		Data:      make(map[string]interface{}),
		mutex:     sync.RWMutex{},
	}
}

func (s *BaseSession) Lock() {
	s.mutex.Lock()
}

func (s *BaseSession) Unlock() {
	s.mutex.Unlock()
}

func (s *BaseSession) RLock() {
	s.mutex.RLock()
}

func (s *BaseSession) RUnlock() {
	s.mutex.RUnlock()
}

// IsExpired 判断会话是否过期
func (s *BaseSession) IsExpired() bool {
	return time.Now().After(s.ExpiredAt)
}

func (s *BaseSession) Get(key string) (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.Data[key]
	return val, ok
}

func (s *BaseSession) Set(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data[key] = value
}

func (s *BaseSession) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Data, key)
}

func (s *BaseSession) ExtendExpiration(ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ExpiredAt = time.Now().Add(ttl)
}

// const (
// 	StateUnauthenticated = uint32(0) // 未鉴权
// 	StateAuthenticated   = uint32(1) // 已鉴权
// 	StateClosed          = uint32(2) // 已关闭
// )

// // Connection 通用连接接口
// type Connection interface {
// 	// 读取消息
// 	ReadMessage() ([]byte, error)

// 	// 发送消息
// 	WriteMessage(message []byte) error

// 	// 关闭连接
// 	Close() error

// 	// 获取远程地址
// 	RemoteAddr() net.Addr

// 	// 设置读取超时时间
// 	SetReadDeadline(t time.Time) error

// 	// 设置写入超时时间
// 	SetWriteDeadline(t time.Time) error
// }

// type Session struct {
// 	ConnId            string
// 	UserId            string
// 	DeviceId          string
// 	Conn              Connection
// 	SendChan          chan []byte
// 	LastHeartbeatTime int64
// 	IP                string
// 	State             atomic.Uint32
// }

// func NewSession(connId string, conn Connection, ip string) *Session {
// 	return &Session{
// 		ConnId:            connId,
// 		Conn:              conn,
// 		SendChan:          make(chan []byte, 1024),
// 		LastHeartbeatTime: time.Now().Unix(),
// 		IP:                ip,
// 		State:             atomic.Uint32{},
// 	}
// }

// // 通过认证后设置用户Id和设备Id
// func (s *Session) SetUserIdAndDeviceId(userId, deviceId string) {
// 	s.UserId = userId
// 	s.DeviceId = deviceId
// 	s.State.Store(StateAuthenticated)
// }

// func (s *Session) IsAuthenticated() bool {
// 	return s.State.Load() == StateAuthenticated
// }

// func (s *Session) SetState(state uint32) {
// 	s.State.Store(state)
// }

// func UpdateHeartbeat(s *Session) {
// 	s.LastHeartbeatTime = time.Now().Unix()
// }

// func (s *Session) IsClosed() bool {
// 	return s.State.Load() == StateClosed
// }

// func (s *Session) Close() {
// 	if s.State.CompareAndSwap(StateAuthenticated, StateClosed) ||
// 		s.State.CompareAndSwap(StateUnauthenticated, StateClosed) {
// 		close(s.SendChan)
// 		_ = s.Conn.Close()
// 	}
// }
