package http_server

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/status_server/internal/handler"
	"github.com/mxxmstar/learning/status_server/status_config"
)

type HttpServer struct {
	engine *gin.Engine
	config *status_config.Config
}

func NewHttpServer(config *status_config.Config) *HttpServer {
	if config.ServerConfig.GlobalConfig.Env == "test" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &HttpServer{
		engine: gin.Default(),
		config: config,
	}

	// 设置中间件
	s.setupMiddleware()

	// 设置路由
	s.setupRoutes()

	return s
}

func (s *HttpServer) setupMiddleware() {
	// 日志中间件
	s.engine.Use(gin.Logger())
	// 错误处理中间件
	s.engine.Use(gin.Recovery())

	// 跨域中间件
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"}, // 允许的请求头
		ExposeHeaders:    []string{"x-jwt-token"},                                       // 允许前端拿到 x-jwt-token 字段，必须要加
		AllowCredentials: true,                                                          // 允许浏览器发送cookie
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境下，允许所有来源
				return true
			}
			// 生产环境下，只允许指定域名
			return strings.Contains(origin, "com.example.com")
		},
		MaxAge: 12 * time.Hour, // 缓存12小时
	}))
}

func (s *HttpServer) setupRoutes() {
	RegisterServerRoutes(s.engine)
}

func (s *HttpServer) Run(c *status_config.Config) error {
	return s.engine.Run(c.GetStatusServerHttpAddr())
}

// 启动定时清理过期服务的任务
func (s *HttpServer) StartCleanupTask() {
	go func() {
		tickerr := time.NewTicker(1 * time.Minute) // 每分钟执行一次
		defer tickerr.Stop()

		for {
			select {
			case <-tickerr.C:
				// 清理过期服务
				// ...
				handler.CleanupExpiredServices()
			}
		}
	}()
}
