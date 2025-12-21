package token_session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	jwt_manager "github.com/mxxmstar/learning/pkg/jwt"
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
	JWTManager  *jwt_manager.JWT
}

// 生成传统的 session 字符串
func GenerateToken(byteCnt int) (string, error) {
	bytes := make([]byte, byteCnt)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewLoginTokenSession(userID uint64, deviceID string, permissions []string, ttl time.Duration) (*LoginTokenSession, error) {
	byteCnt := 16
	token, err := GenerateToken(byteCnt)
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

// 创建支持 JWT 的登录令牌会话
func NewLoginTokenSessionWithJWT(userID uint64, deviceID string, permissions []string, ttl time.Duration, jwtManager *jwt_manager.JWT) (*LoginTokenSession, error) {
	session, err := NewLoginTokenSession(userID, deviceID, permissions, ttl)
	if err != nil {
		return nil, err
	}

	// 设置JWT管理器
	session.JWTManager = jwtManager

	return session, nil
}

func (lts *LoginTokenSession) GenerateJWTToken() (string, error) {
	lts.BaseSession.RLock()
	defer lts.BaseSession.RUnlock()

	if lts.JWTManager == nil {
		return "", session.ErrJWTManagerNotSet
	}

	return lts.JWTManager.GenerateToken(lts.UserID, lts.DeviceID)
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
