package service

import (
	"context"
	"errors"

	"github.com/mxxmstar/learning/verify_server/internal/domain"
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

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func validateUser(u domain.User) error {
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

// 注册新用户
func (s *UserService) Signup(ctx context.Context, u domain.User) error {
	if err := validateUser(u); err != nil {
		return err
	}
	err := s.repo.CreateUser(ctx, &domain.User{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password, // 先不加密，后续再处理
		CTime:    u.CTime,
	})

}

// // 更新用户信息
// func (s *UserService) UpdateUser(ctx context.Context, userId string, updateInfo map[string]interface{}) error {

// }

// // 删除用户
// func (s *UserService) DeleteUser(ctx context.Context, userId string) error {

// }
