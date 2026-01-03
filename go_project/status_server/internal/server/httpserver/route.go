package http_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/pkg/store/redis"
	"github.com/mxxmstar/learning/status_server/status_config"
)

// 注册查询服务路由
func RegisterQueryRoutes(server *gin.Engine, cfg *status_config.Config, redisClient *redis.RedisClient) {

}

// 注册服务管理路由
func RegisterServerRoutes(server *gin.Engine, cfg *status_config.Config, redisClient *redis.RedisClient) {
	// 服务注册相关路由
	api := server.Group("/api")
	{
		// 服务注册相关路由
		api.POST("/service/register", handler.ServiceRegisterHandler)
		api.POST("/service/deregister", handler.ServiceDeregisterHandler)
		api.POST("/service/heartbeat", handler.ServiceHeartbeatHandler)

		// 服务发现相关路由
		api.POST("/discovery/by-name", handler.ServiceDiscoveryByNameHandler)
		api.POST("/discovery/by-tags", handler.ServiceDiscoveryByTagsHandler)
		api.POST("/discovery/by-metadata", handler.ServiceDiscoveryByMetadataHandler)

		// 服务状态查询路由
		api.POST("/status/query", handler.ServiceStatusQueryHandler)
		api.GET("/status/overview", handler.ServiceStatusOverviewHandler)

		// 健康检查
		api.GET("/health", handler.HealthCheckHandler)
	}

	// 为 Gate Server 提供的特定路由
	gate := server.Group("/gate")
	{
		gate.POST("/discovery/verify", handler.ServiceDiscoveryByTagsHandler)
		gate.POST("/service/register", handler.ServiceRegisterHandler)
		gate.POST("/service/heartbeat", handler.ServiceHeartbeatHandler)
		gate.POST("/discovery/by-tags", handler.ServiceDiscoveryByTagsHandler)
	}

	// 根路径
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Status Server is running",
			"version": "1.0.0",
			"status":  "healthy",
		})
	})
}
