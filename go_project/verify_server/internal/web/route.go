package web

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/pkg/database"
	"github.com/mxxmstar/learning/pkg/store/redis"
	"github.com/mxxmstar/learning/verify_server/internal/repository"
	"github.com/mxxmstar/learning/verify_server/internal/repository/dao"
	"github.com/mxxmstar/learning/verify_server/internal/service"
	"github.com/mxxmstar/learning/verify_server/internal/web/handler"
	"github.com/mxxmstar/learning/verify_server/verify_config"
)

func RegisterUserRoutes(server *gin.Engine, cfg *verify_config.Config, db database.DBInterface, redisClient *redis.RedisClient) {

	// 初始化仓库
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)

	// 初始化服务
	authService := service.NewAuthService(userRepo, redisClient, cfg.VerifyService.JWTSecret, cfg.VerifyService.TokenLifeTime)
	userService := service.NewUserService(userRepo)

	// 注册用户验证处理器
	authHandler := handler.NewAuthHandler(authService, userService)
	// 注册用户处理器
	userHandler := handler.NewUserHandler(userService)

	log.Printf("============%s", cfg.ServerConfig.GlobalConfig.Env)
	// 注册用户注册相关路由（测试用）
	if cfg.ServerConfig.GlobalConfig.Env == "test" {
		authGroup := server.Group("/user-auth")
		{
			authGroup.POST("/signup", authHandler.SignupHandler)
			authGroup.POST("/login", authHandler.LoginHandler)
			authGroup.POST("/oauth", authHandler.OAuthHandler)
		}
	}

	// 注册用户注册相关路由（与 gate 通信）
	gateAuthGroup := server.Group("gate/user-auth")
	{
		gateAuthGroup.POST("/verify-session", authHandler.VerifySessionHandler)
		gateAuthGroup.POST("/verify-jwt", authHandler.VerifyJWTHandler)
		gateAuthGroup.POST("/refresh-session", authHandler.RefreshSessionHandler)
	}

	// 注册用户相关路由（测试用）
	if cfg.ServerConfig.GlobalConfig.Env == "test" {
		userGroup := server.Group("/user")
		{
			userGroup.GET("/profile", userHandler.ProfileHandler)
			userGroup.PUT("/profile", userHandler.UpdateProfileHandler)
		}
	}
	gateUserGroup := server.Group("gate/user")
	{
		gateUserGroup.GET("/profile", userHandler.ProfileHandler)
		gateUserGroup.PUT("/profile", userHandler.UpdateProfileHandler)
	}

}
