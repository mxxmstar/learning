package web

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/verify_server/internal/repository"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
	"github.com/mxxmstar/learning/verify_server/internal/service"
	"github.com/mxxmstar/learning/verify_server/internal/web/handler"
	"github.com/mxxmstar/learning/verify_server/verify_config"
)

func RegisterUserRoutes(server *gin.Engine, cfg *verify_config.Config) {
	// 初始化数据库
	db, err := verify_config.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db, cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Database tables initialized successfully.")

	// 初始化Redis客户端
	redisClient, err := verify_config.InitRedis(cfg)
	if err != nil {
		panic(err)
	}

	// 初始化仓库
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)

	// 初始化服务
	authService := service.NewAuthService(userRepo, redisClient)
	userService := service.NewUserService(userRepo)

	// 注册用户处理器
	authHandler := handler.NewAuthHandler(authService, userService)
	// 注册用户处理器
	userHandler := handler.NewUserHandler(userService)

	// 注册用户注册相关路由
	authGroup := server.Group("/user/auth")
	{
		authGroup.POST("/signup", authHandler.SignupHandler)
		authGroup.POST("/login", authHandler.LoginHandler)
		authGroup.POST("/oauth", authHandler.OAuthHandler)
	}
	// 注册用户相关路由
	userGroup := server.Group("/user")
	{
		userGroup.GET("/profile", userHandler.ProfileHandler)
		userGroup.PUT("/profile", userHandler.UpdateProfileHandler)
	}

}
