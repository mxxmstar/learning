package database

import (
	"context"

	"gorm.io/gorm"
)

// Wrapper 包装器，用于封装不同的DAO和数据库操作

// 抽象数据库上下文接口
type DBContextInterface interface {
	// 定义数据库上下文操作方法的接口
	Create(value interface{}) DBContextInterface
	Where(query interface{}, args ...interface{}) DBContextInterface
	First(dest interface{}) DBContextInterface
	Error() error
}

// 描述业务需要的数据库操作方法
// 抽象数据库接口，避免直接依赖具体的ORM实现
type DBInterface interface {
	// 在这里定义数据库操作方法的接口
	WithContext(ctx context.Context) DBContextInterface
}

// GORMContextWrapper 包装*gorm.DB实现DBContextInterface
type GORMContextWrapper struct {
	db *gorm.DB
}

// GORM具体适配
// GORMWrapper 包装gorm.DB以实现DBInterface
type GORMWrapper struct {
	db *gorm.DB
}

func NewGORMWrapper(db *gorm.DB) *GORMWrapper {
	return &GORMWrapper{db: db}
}

func (w *GORMWrapper) WithContext(ctx context.Context) DBContextInterface {
	return &GORMContextWrapper{db: w.db.WithContext(ctx)}
}

// Create 实现DBContextInterface的Create方法
func (w *GORMContextWrapper) Create(value interface{}) DBContextInterface {
	w.db = w.db.Create(value).Statement.DB
	return w
}

func (w *GORMContextWrapper) Where(query interface{}, args ...interface{}) DBContextInterface {
	return &GORMContextWrapper{db: w.db.Where(query, args...)}
}

func (w *GORMContextWrapper) First(dest interface{}) DBContextInterface {
	return &GORMContextWrapper{db: w.db.First(dest)}
}

func (w *GORMContextWrapper) Error() error {
	return w.db.Error
}
