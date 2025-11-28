package service

import (
	"context"
	"errors"

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
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Signup(ctx context.Context, user *domain.User) error {
	// 在这里调用 encrypt 对密码进行加密

	return s.userRepo.CreateUser(ctx, user)
}
