package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	jwt_manager "github.com/mxxmstar/learning/pkg/jwt"
	"github.com/mxxmstar/learning/pkg/session/token_session"
	"github.com/mxxmstar/learning/pkg/store/redis"
	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/repository"
)

var (
	// ErrInvalidCredentials 表示用户名或密码错误
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrUserEmailConflict 表示邮箱冲突错误
	ErrUserEmailConflict = repository.ErrDuplicateEmail

	// ErrUserUsernameConflict 表示用户名冲突错误
	ErrUserUsernameConflict = repository.ErrDuplicateUsername

	// ErrUserDisabled 表示用户账户已被禁用
	ErrUserDisabled = errors.New("user account is disabled")

	// ErrTooManyLoginAttempts 表示登录尝试次数过多
	ErrTooManyLoginAttempts = errors.New("too many login attempts")

	// ErrInvalidUserInfo 表示用户提交的信息不符合要求
	ErrInvalidUserInfo = errors.New("invalid user information")

	// SessionTTl 表示会话过期时间，默认24小时
	SessionTTL = 24 * time.Hour
)

type AuthService struct {
	userRepo      *repository.UserRepository
	redisClient   *redis.RedisClient
	jwtSecret     string
	tokenLifeTime int
	jwtManager    *jwt_manager.JWT
}

func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.RedisClient, jwtSecret string, tokenLifetime int) *AuthService {
	if tokenLifetime <= 0 {
		// 默认设置为1小时
		tokenLifetime = 3600
	}

	jwtMgr := jwt_manager.NewJWT([]byte(jwtSecret), "verify_server", tokenLifetime)
	return &AuthService{
		userRepo:      userRepo,
		redisClient:   redisClient,
		jwtSecret:     jwtSecret,
		tokenLifeTime: tokenLifetime,
		jwtManager:    jwtMgr,
	}
}

func (s *AuthService) Signup(ctx context.Context, user *domain.User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	// 在这里调用 encrypt 对密码进行加密

	return s.userRepo.CreateUser(ctx, user)
}

// Login 用户登录并创建session，这里不通过 session 对象生成 JWT，在handler层统一整合
func (s *AuthService) LoginByEmail(ctx context.Context, email, password string, loginCtx *domain.LoginContext) (string, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// 验证密码（这里应该使用加密验证）
	if user.Password != password {
		return "", ErrInvalidCredentials
	}

	// 创建session
	// TODO: 权限待完善
	deviceID := ""
	if loginCtx != nil {
		deviceID = loginCtx.DeviceId
	}
	loginSession, err := token_session.NewLoginTokenSession(
		user.Id,
		deviceID,   // 设备ID
		[]string{}, // 权限列表
		SessionTTL,
	)
	if err != nil {
		return "", err
	}

	// 序列化用户信息
	userData, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// 存储 session 到 Redis
	t := loginSession.GetToken()
	key := "session:" + t
	err = s.redisClient.Set(ctx, key, string(userData), SessionTTL)
	if err != nil {
		return "", err
	}

	// TODO: loginCtx信息的管理

	return t, nil
}

// GenerateJWT 生成JWT令牌
func (s *AuthService) GenerateJWT(user *domain.User, loginCtx *domain.LoginContext) (string, error) {
	userID := user.Id
	deviceID := ""
	if loginCtx != nil {
		deviceID = loginCtx.DeviceId
	}

	// 生成JWT token
	token, err := s.jwtManager.GenerateToken(userID, deviceID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetSessionUser 从session中获取用户信息
func (s *AuthService) GetSessionUser(ctx context.Context, sessionID string) (*domain.User, error) {
	key := "session:" + sessionID
	userStr, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Logout 用户登出，清除session
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	return s.redisClient.Del(ctx, key)
}

// RefreshSession 刷新session的过期时间
func (s *AuthService) RefreshSession(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	return s.redisClient.Expire(ctx, key, SessionTTL)
}

// 验证并解析 JWT 令牌
func (s *AuthService) ValidateAndParseJWT(token string) (*jwt_manager.CustomClaims, error) {
	// 解析JWT
	claims, err := s.jwtManager.ParseToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func validateUser(u *domain.User) error {
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return ErrInvalidUserInfo
	}
	// 校验邮箱格式
	// if !isValidEmail(u.Email) {
	// 	return ErrInvalidUserInfo
	// }
	if len(u.Password) < 6 {
		return ErrInvalidUserInfo
	}
	return nil
}
