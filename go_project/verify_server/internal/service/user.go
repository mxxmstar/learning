package service

import (
	"errors"

	"github.com/mxxmstar/learning/verify_server/internal/repository"
)

var (
	// ErrDuplicateEmail 表示邮箱已经被注册使用
	ErrDuplicateEmail = errors.New("email already registered")

	// ErrDuplicateUsername 表示用户名已经被注册使用
	ErrDuplicateUsername = errors.New("username already registered")

	// ErrInvalidUserInfo 表示用户提交的信息不符合要求
	ErrInvalidUserInfo = errors.New("invalid user information")

	// ErrUserNotFound 表示未找到指定用户
	ErrUserNotFound = errors.New("user not found")
)

// 调用领域对象的方法，实现业务逻辑
type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// // 更新用户信息
// func (s *UserService) UpdateUser(ctx context.Context, userId string, updateInfo map[string]interface{}) error {

// }

// // 删除用户
// func (s *UserService) DeleteUser(ctx context.Context, userId string) error {

// }
