package repository

import (
	"context"

	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
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
