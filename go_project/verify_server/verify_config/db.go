package verify_config

import (
	"fmt"

	"github.com/mxxmstar/learning/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) (database.DBInterface, error) {
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

	println("InitDB")

	g, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	println("gorm open success")
	sqlDB, err := g.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(cfg.Database.Pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.Pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.Pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.Pool.ConnMaxIdleTime)

	db := database.NewGORMWrapper(g)
	return db, nil
}
