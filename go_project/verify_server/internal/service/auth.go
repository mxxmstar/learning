package service

import (
	"context"
	"errors"
	"time"

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

	// SessionTTl 表示会话过期时间，默认24小时
	SessionTTl = 24 * time.Hour
)

type AuthService struct {
	userRepo    *repository.UserRepository
	redisClient *redis.RedisClient
}

func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.RedisClient) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *AuthService) Signup(ctx context.Context, user *domain.User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	// 在这里调用 encrypt 对密码进行加密

	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) LoginByEmail(ctx context.Context, email, password string) (*domain.User, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return &domain.User{}, err
	}

	// 比较密码

	// 存储session

	return user, nil
}

// Login 用户登录并创建session
// func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
// 	// 查找用户
// 	user, err := s.userRepo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return "", ErrInvalidCredentials
// 	}

// 	// 验证密码（这里应该使用加密验证）
// 	if user.Password != password {
// 		return "", ErrInvalidCredentials
// 	}

// 	// 创建session ID
// 	sessionID := generateSessionID()

// 	// 序列化用户信息
// 	userData, err := json.Marshal(user)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 存储到Redis
// 	key := "session:" + sessionID
// 	err = s.redisClient.Set(ctx, key, string(userData), SessionTTL)
// 	if err != nil {
// 		return "", err
// 	}

// 	return sessionID, nil
// }

// // GetSessionUser 从session中获取用户信息
// func (s *AuthService) GetSessionUser(ctx context.Context, sessionID string) (*domain.User, error) {
// 	key := "session:" + sessionID
// 	userStr, err := s.redisClient.Get(ctx, key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var user domain.User
// 	err = json.Unmarshal([]byte(userStr), &user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// // Logout 用户登出
// func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
// 	key := "session:" + sessionID
// 	return s.redisClient.Del(ctx, key)
// }

// // generateSessionID 生成session ID（简化版，实际应使用更安全的方式）
// func generateSessionID() string {
// 	// 实际应用中应该使用更安全的随机字符串生成方法
// 	bytes := make([]byte, 16)
// 	// 此处省略实际的随机数生成逻辑
// 	return "session_" + string(bytes)
// }

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
