package session

import (
	"net"
	"sync/atomic"
	"time"
)

const (
	StateUnauthenticated = uint32(0) // 未鉴权
	StateAuthenticated   = uint32(1) // 已鉴权
	StateClosed          = uint32(2) // 已关闭
)

// Connection 通用连接接口
type Connection interface {
	// 读取消息
	ReadMessage() ([]byte, error)

	// 发送消息
	WriteMessage(message []byte) error

	// 关闭连接
	Close() error

	// 获取远程地址
	RemoteAddr() net.Addr

	// 设置读取超时时间
	SetReadDeadline(t time.Time) error

	// 设置写入超时时间
	SetWriteDeadline(t time.Time) error
}

type Session struct {
	ConnID            string
	UserID            string
	DeviceID          string
	Conn              Connection
	SendChan          chan []byte
	LastHeartbeatTime int64
	IP                string
	State             atomic.Uint32
}

func NewSession(connID string, conn Connection, ip string) *Session {
	return &Session{
		ConnID:            connID,
		Conn:              conn,
		SendChan:          make(chan []byte, 1024),
		LastHeartbeatTime: time.Now().Unix(),
		IP:                ip,
		State:             atomic.Uint32{},
	}
}

// 通过认证后设置用户ID和设备ID
func (s *Session) SetUserIDAndDeviceID(userID, deviceID string) {
	s.UserID = userID
	s.DeviceID = deviceID
	s.State.Store(StateAuthenticated)
}

func (s *Session) IsAuthenticated() bool {
	return s.State.Load() == StateAuthenticated
}

func (s *Session) SetState(state uint32) {
	s.State.Store(state)
}

func UpdateHeartbeat(s *Session) {
	s.LastHeartbeatTime = time.Now().Unix()
}

func (s *Session) IsClosed() bool {
	return s.State.Load() == StateClosed
}

func (s *Session) Close() {
	if s.State.CompareAndSwap(StateAuthenticated, StateClosed) ||
		s.State.CompareAndSwap(StateUnauthenticated, StateClosed) {
		close(s.SendChan)
		_ = s.Conn.Close()
	}
}
