package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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
}

func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.RedisClient, jwtSecret string, tokenLifetime int) *AuthService {
	if tokenLifetime <= 0 {
		// 默认设置为1小时
		tokenLifetime = 3600
	}
	return &AuthService{
		userRepo:      userRepo,
		redisClient:   redisClient,
		jwtSecret:     jwtSecret,
		tokenLifeTime: tokenLifetime,
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
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
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
	// TODO: 设备ID和权限待完善
	loginSession, err := token_session.NewLoginTokenSession(
		user.Id,
		"",         // 设备ID
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

	// 存储到Redis
	t := loginSession.GetToken()
	key := "session:" + t
	err = s.redisClient.Set(ctx, key, string(userData), SessionTTL)
	if err != nil {
		return "", err
	}

	return t, nil
}

// func (s *AuthService) GenerateJWT(user *domain.User) (string, error) {
// 	// 创建声明
// 	claims := jwt.MapClaims{
// 		"user_id":   user.Id,
// 		"email":     user.Email,
// 		"username":  user.Username,
// 		"exp":       time.Now().Add(time.Second * time.Duration(s.tokenLifeTime)).Unix(),
// 		"issued_at": time.Now().Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(s.jwtSecret))
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
// 	// 这里暂时保留原逻辑，但实际应该被替换
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
