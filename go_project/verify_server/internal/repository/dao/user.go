package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = gorm.ErrRecordNotFound
	// 邮箱冲突
	ErrEmailConflict = errors.New("email already exists")
)

type UserDAO struct {
	db *gorm.DB
	// cache Cache      // 缓存
	// logger Logger     // 日志
	// timeout time.Duration // 超时配置
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Email        string `gorm:"size:128;uniqueIndex;not null"`
	Phone        string `gorm:"size:32"`
	AvatarURL    string `gorm:"size:255"`

	IsBanned    bool       `gorm:"default:false"`
	LastLoginAt *time.Time `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (dao *UserDAO) Insert(ctx context.Context, user *User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}
