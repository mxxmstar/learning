package repository

import (
	"context"

	"github.com/mxxmstar/learning/pkg/database"
	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
)

// 将 error_converter 中的错误转换为 repository 业务层的错误
var (
	// ErrDuplicateEmail 表示邮箱冲突错误
	ErrDuplicateEmail = database.ErrEmailConflict

	// ErrUserNotFound 表示用户不存在错误
	ErrUserNotFound = database.ErrUserNotFound
)

type UserRepository struct {
	userDAO *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{
		userDAO: userDAO,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return repo.userDAO.Insert(ctx, &dao.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
}
