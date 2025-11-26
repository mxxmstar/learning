package web

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/verify_server/verify_config"
)

func RegisterRoutes(server *gin.Engine, cfg *verify_config.Config) {
	// 初始化服务
	// 初始化处理器
	// 注册中间件
	RegisterUserRoutes(server, cfg)

}

func InitWebServer(cfg *verify_config.Config) *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:3000"},
		// 不写就默认所有请求
		// AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 允许前端拿到 x-jwt-token 字段，必须要加
		ExposeHeaders: []string{"x-jwt-token"},
		// 允许浏览器发送cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境下，允许所有来源
				return true
			}
			// 生产环境下，只允许指定域名
			return strings.Contains(origin, "com.example.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}
