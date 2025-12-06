package gatewaysession

import (
	"time"

	"github.com/mxxmstar/learning/pkg/session"
)

type ClientGatewaySession struct {
	*session.BaseSession
	Connection     interface{} // tcp或者websocket 连接
	ClientIP       string
	UserID         uint64
	DeviceID       string
	LastActiveTime time.Time
	HeartbeatCount int64
}

func NewClientGatewaySession(sessionID session.SessionID, conn interface{}, clientIP string, ttl time.Duration) *ClientGatewaySession {
	// 客户端刚连接到gateway时，还未进行身份验证，UserID和DeviceID为空
	return &ClientGatewaySession{
		BaseSession:    session.NewBaseSession(sessionID, session.ClientGatewaySession, ttl),
		Connection:     conn,
		ClientIP:       clientIP,
		LastActiveTime: time.Now(),
	}
}

// SetUser 设置用户ID和DeviceID
func (s *ClientGatewaySession) SetUser(userID uint64, deviceID string) {
	s.BaseSession.Lock()
	defer s.BaseSession.Unlock()
	s.UserID = userID
	s.DeviceID = deviceID
	// 存储用户ID和DeviceID到data中
	s.BaseSession.Set("user_id", userID)
	s.BaseSession.Set("device_id", deviceID)
}

// GetUser 获取用户ID和DeviceID
func (s *ClientGatewaySession) GetUser() (uint64, string) {
	s.BaseSession.RLock()
	defer s.BaseSession.RUnlock()
	return s.UserID, s.DeviceID
}

// UpdateLastActiveTime 更新最后活跃时间
func (s *ClientGatewaySession) UpdateLastActiveTime() {
	s.BaseSession.Lock()
	defer s.BaseSession.Unlock()
	s.LastActiveTime = time.Now()
	s.HeartbeatCount++
}

func (s *ClientGatewaySession) GetConnection() interface{} {
	s.BaseSession.RLock()
	defer s.BaseSession.RUnlock()
	return s.Connection
}

func (s *ClientGatewaySession) GetClientIP() string {
	s.BaseSession.RLock()
	defer s.BaseSession.RUnlock()
	return s.ClientIP
}
