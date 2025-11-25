package service

import (
	"errors"

	"github.com/mxxmstar/learning/verify_server/internal/repository"
)

var (
	// ErrInvalidCredentials 表示用户名或密码错误
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrUserDisabled 表示用户账户已被禁用
	ErrUserDisabled = errors.New("user account is disabled")

	// ErrTooManyLoginAttempts 表示登录尝试次数过多
	ErrTooManyLoginAttempts = errors.New("too many login attempts")
)

type AuthService struct {
	userRepo *repository.UserRepository
}
