package token_session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/mxxmstar/learning/pkg/session"
)

// LoginTokenSession 登录令牌会话
type LoginTokenSession struct {
	*session.BaseSession
	Token       string
	UserID      uint64
	DeviceID    string
	CreatedAt   time.Time
	Permissions []string // 用户权限列表
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewLoginTokenSession(userID uint64, deviceID string, permissions []string, ttl time.Duration) (*LoginTokenSession, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	session := &LoginTokenSession{
		BaseSession: session.NewBaseSession(session.SessionID("token_"+token), session.LoginTokenSession, ttl),
		Token:       token,
		UserID:      userID,
		DeviceID:    deviceID,
		CreatedAt:   time.Now(),
		Permissions: permissions,
	}

	// 存储关键信息到Data
	session.Data["user_id"] = userID
	session.Data["device_id"] = deviceID
	session.Data["permissions"] = permissions

	return session, nil
}

func (lts *LoginTokenSession) GetToken() string {
	lts.BaseSession.RLock()
	defer lts.BaseSession.RUnlock()
	return lts.Token
}

func (lts *LoginTokenSession) ValidatePermission(permission string) bool {
	lts.BaseSession.RLock()
	defer lts.BaseSession.RUnlock()
	for _, p := range lts.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (lts *LoginTokenSession) GetPermissions() []string {
	lts.BaseSession.RLock()
	defer lts.BaseSession.RUnlock()
	return lts.Permissions
}
