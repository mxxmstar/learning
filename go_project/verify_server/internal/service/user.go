package service

import (
	"errors"

	"github.com/mxxmstar/learning/verify_server/internal/repository"
)

var (
	// ErrDuplicateEmail 表示邮箱已经被注册使用
	ErrDuplicateEmail = errors.New("email already registered")

	// ErrInvalidUserInfo 表示用户提交的信息不符合要求
	ErrInvalidUserInfo = errors.New("invalid user information")

	// ErrUserNotFound 表示未找到指定用户
	ErrUserNotFound = errors.New("user not found")
)

// 调用领域对象的方法，实现业务逻辑
type UserService struct {
	repo *repository.UserRepository
}
