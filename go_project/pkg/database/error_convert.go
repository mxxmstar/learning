package database

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	// 邮箱冲突
	ErrEmailConflict = errors.New("email already exists")
)

// DBErrorConverter 数据库错误转换器接口
type DBErrorConverter interface {
	ConvertError(err error) error
}

// GORMErrorConverter GORM错误转换器
type GORMErrorConverter struct{}

func (c *GORMErrorConverter) ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// 处理GORM特定错误
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrEmailConflict
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrUserNotFound
	default:
		return err
	}
}
