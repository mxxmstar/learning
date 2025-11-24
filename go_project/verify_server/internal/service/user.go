package service

import "github.com/mxxmstar/learning/verify_server/internal/repository"

var ErrDuplicatedKey = repository.ErrDuplicatedKey
var ErrInvalidUserOrPassword = repository.ErrInvalidUserOrPassword

// 调用领域对象的方法，实现业务逻辑
type UserService struct {
	repo *repository.UserRepository
}
