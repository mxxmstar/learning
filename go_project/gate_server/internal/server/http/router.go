package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/gate_server/gate_config"
	"github.com/mxxmstar/learning/verify_server/internal/web/handler"
)

func RegisterUserRoutes(server *gin.Engine, cfg *gate_config.Config) {

	// 注册用户验证处理器
	authHandler := handler.NewAuthHandler(authService, userService)
	// 注册用户处理器
	userHandler := handler.NewUserHandler(userService)

	// 注册用户注册相关路由
	authGroup := server.Group("/user-auth")
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
