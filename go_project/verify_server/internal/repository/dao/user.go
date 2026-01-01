package dao

import (
	"context"
	"time"

	"github.com/mxxmstar/learning/pkg/database"
)

type UserDAO struct {
	db             database.DBInterface      // 数据库接口
	errorConverter database.DBErrorConverter // 数据库错误转换器
	// cache Cache      // 缓存
	// logger Logger     // 日志
	// timeout time.Duration // 超时配置
}

func NewUserDAO(db database.DBInterface) *UserDAO {
	return &UserDAO{
		db:             db,
		errorConverter: &database.GORMErrorConverter{},
	}
}

// SetErrorConverter 允许外部设置自定义错误转换器
func (dao *UserDAO) SetErrorConverter(converter database.DBErrorConverter) {
	dao.errorConverter = converter
}

type User struct {
	// 用户唯一主键Id 自动递增
	Id uint64 `gorm:"primaryKey;autoIncrement"`
	// 用户名 64字节 唯一索引 不能为空
	Username string `gorm:"size:64;uniqueIndex;not null"`
	// 密码 255字节 不能为空
	Password string `gorm:"size:255;not null"`
	// 邮箱 128字节 唯一索引 不能为空
	Email string `gorm:"size:128;uniqueIndex;not null"`
	// 手机号 32字节 唯一索引
	Phone string `gorm:"size:32"`
	// 头像URL 255字节 为空时表示使用默认图像
	AvatarURL string `gorm:"size:255"`

	// 是否被封禁 默认值为 false
	IsBanned bool `gorm:"default:false"`
	// 最后登录的时间戳 索引 0表示从未登录
	LastLoginAt int64 `gorm:"index"`
	// 记录创建和更新时间 自动管理
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (dao *UserDAO) Insert(ctx context.Context, user *User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastLoginAt = 0
	dbCtx := dao.db.WithContext(ctx).Create(user)
	if dbCtx.Error() != nil {
		return dao.errorConverter.ConvertError(dbCtx.Error())
	}
	return nil
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	dbCtx := dao.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if dbCtx.Error() != nil {
		return nil, dao.errorConverter.ConvertError(dbCtx.Error())
	}
	return &user, nil
}
