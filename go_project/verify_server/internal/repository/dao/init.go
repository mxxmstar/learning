package dao

import (
	"github.com/mxxmstar/learning/pkg/database"
	"github.com/mxxmstar/learning/verify_server/verify_config"
)

// InitTables 初始化数据库表, 仅在非生产环境下自动迁移
func InitTables(db database.DBInterface, cfg *verify_config.Config) error {
	// if cfg.Database.AutoMigrate && cfg.Env != "production" {
	// 	return db.AutoMigrate(&User{})
	// }
	return db.AutoMigrate(&User{})
	// return nil
}
