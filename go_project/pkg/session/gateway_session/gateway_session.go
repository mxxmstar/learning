package gatewaysession

import (
	"time"

	"github.com/mxxmstar/learning/pkg/session"
)

type ClientGatewaySession struct {
	*session.BaseSession
	Connection     interface{} // tcp或者websocket 连接
	ClientIP       string
	UserId         uint64
	DeviceId       string
	LastActiveTime time.Time
	HeartbeatCount int64
}

func NewClientGatewaySession(sessionId session.SessionId, conn interface{}, clientIP string, ttl time.Duration) *ClientGatewaySession {
	// 客户端刚连接到gateway时，还未进行身份验证，UserId和DeviceId为空
	return &ClientGatewaySession{
		BaseSession:    session.NewBaseSession(sessionId, session.ClientGatewaySession, ttl),
		Connection:     conn,
		ClientIP:       clientIP,
		LastActiveTime: time.Now(),
	}
}

// SetUser 设置用户Id和DeviceId
func (s *ClientGatewaySession) SetUser(userId uint64, deviceId string) {
	s.BaseSession.Lock()
	defer s.BaseSession.Unlock()
	s.UserId = userId
	s.DeviceId = deviceId
	// 存储用户Id和DeviceId到data中
	s.BaseSession.Set("user_id", userId)
	s.BaseSession.Set("device_id", deviceId)
}

// GetUser 获取用户Id和DeviceId
func (s *ClientGatewaySession) GetUser() (uint64, string) {
	s.BaseSession.RLock()
	defer s.BaseSession.RUnlock()
	return s.UserId, s.DeviceId
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
