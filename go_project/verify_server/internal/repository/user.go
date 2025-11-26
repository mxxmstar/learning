package repository

import (
	"context"
	"errors"

	"github.com/mxxmstar/learning/pkg/database"
	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
)

// 将 error_converter 中的错误转换为 repository 业务层的错误
var (
	// ErrUserNotFound 表示用户不存在错误
	ErrUserNotFound = database.ErrUserNotFound

	// ErrDuplicateEmail 表示邮箱冲突错误
	ErrDuplicateEmail = database.ErrEmailConflict

	// ErrDuplicateUsername 表示用户名冲突错误
	ErrDuplicateUsername = database.ErrUsernameConflict
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
	err := repo.userDAO.Insert(ctx, &dao.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if errors.Is(err, database.ErrEmailConflict) {
		return ErrDuplicateEmail
	}
	if errors.Is(err, database.ErrUsernameConflict) {
		return ErrDuplicateUsername
	}
	return err
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := repo.userDAO.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &domain.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
