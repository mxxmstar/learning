package web

import (
	"fmt"

	regexp "github.com/dlclark/regexp2"
	"github.com/mxxmstar/learning/verify_server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserHandler struct {
	// sver

	// emailExp 邮箱正则表达式
	emailExp *regexp.Regexp
	// passwordExp 密码正则表达式
	passwordExp *regexp.Regexp
}

type DB interface {
	AutoMigrate(dst ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dbType := cfg.Database.Type
	dsn := cfg.Database.AuthDB.DSN
	var dialector gorm.Dialector
	fmt.Printf("dbType: %s, dsn: %s\n", dbType, dsn)
	switch dbType {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		// TODO: postgres
		break
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(cfg.Database.Pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.Pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.Pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.Pool.ConnMaxIdleTime)
	return db, nil
}
