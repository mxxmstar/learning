package dao

import (
	"github.com/mxxmstar/learning/pkg/config"
	"gorm.io/gorm"
)

// InitTables 初始化数据库表, 仅在非生产环境下自动迁移
func InitTables(db *gorm.DB, cfg *config.Config) error {
	if cfg.Database.AutoMigrate && cfg.Env != "production" {
		return db.AutoMigrate(&User{})
	}
	return nil
}
