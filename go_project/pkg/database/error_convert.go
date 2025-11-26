package database

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	// 邮箱冲突
	ErrEmailConflict = errors.New("email already exists")
	// 用户名冲突
	ErrUsernameConflict = errors.New("username already exists")
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
		// 提取具体的冲突信息
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "username"):
			return ErrUsernameConflict
		case strings.Contains(errMsg, "email"):
			return ErrEmailConflict
		default:
			// 默认返回用户名冲突错误
			return ErrUsernameConflict
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrUserNotFound
	default:
		return err
	}
}
