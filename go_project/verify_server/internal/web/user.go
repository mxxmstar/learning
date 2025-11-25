package web

import (
	"fmt"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mxxmstar/learning/pkg/database"
	"github.com/mxxmstar/learning/verify_server/config"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
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

func initDB(cfg *config.Config) (database.DBInterface, error) {
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
	g, err := gorm.Open(dialector, &gorm.Config{})
	if g != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

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

func RegisterUserRoutes(server *gin.Engine, cfg *config.Config) {
	db, err := initDB(cfg)
	if err != nil {
		panic(err)
	}
	userDAO := dao.NewUserDAO(db)

}
