package web

import (
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/verify_server/config"
)

func RegisterRoutes(server *gin.Engine, cfg *config.Config) {
	// 初始化服务
	// 初始化处理器
	// 注册中间件

	// 注册路由(不需要认证)
	authGroup := server.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.SignupHandler)
		authGroup.POST("/login", authHandler.SignupHandler)
		authGroup.POST("/oauth/:provider", authHandler.OAuthHandler)
	}

	// 注册路由(需要认证)
	// apiGroup := server.Group("/api")
	// apiGroup.Use(jwtMiddleware)
	// {
	// 	apiGroup.GET("/profile", userHandler.ProfileHandler)
	// 	apiGroup.PUT("/profile", userHandler.UpdateProfileHandler)
	// }
}
