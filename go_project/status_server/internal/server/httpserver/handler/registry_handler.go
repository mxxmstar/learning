package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	status_def "github.com/mxxmstar/learning/pkg/def/status"
	status_model "github.com/mxxmstar/learning/pkg/model/status"
	etcd_registry "github.com/mxxmstar/learning/status_server/internal/registry/etcd_registry"
	mem_registry "github.com/mxxmstar/learning/status_server/internal/registry/mem_registry"
)

const (
	EnableRegistry    = "memory"
	HeartbeatInterval = 60 // TODO: 在配置中配置
)

var registryInstance status_model.Registry

func InitRegistry(endpoints []string) error {
	if EnableRegistry == "memory" {
		registryInstance = mem_registry.NewMemRegistry()
	} else if EnableRegistry == "etcd" {
		reg, err := etcd_registry.NewEtcdRegistry(endpoints)
		if err != nil {
			return err
		}
		registryInstance = reg
	} else if EnableRegistry == "redis" {
		return errors.New("redis registry not supported")
	} else {
		return errors.New("Invalid registry type")
	}

	return nil
}

// ServiceRegisterHandler 服务注册处理器
func ServiceRegisterHandler(c *gin.Context) {
	var req status_def.ServiceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, status_def.ServiceRegisterResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	serviceInfo := req.ConvertToStatusServiceInfo()

}

// ServiceDiscoveryByTagsHandler 服务发现处理器
func ServiceDiscoveryByTagsHandler(c *gin.Context) {
	var req status.ServiceDiscoveryByTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, status.ServiceDiscoveryByTagsResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	services, err := registryInstance.DiscoverServicesByType(req.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, status.ServiceDiscoveryByTagsResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}

	if len(services) == 0 {
		c.JSON(http.StatusOK, status.ServiceDiscoveryByTagsResponse{
			Code:     404,
			Message:  "Service not found",
			Services: nil,
		})
		return
	}

	// 根据策略选择服务（这里简化为选择第一个）
	selectedService := services[0]

	c.JSON(http.StatusOK, status.ServiceDiscoveryByTagsResponse{
		Code:     200,
		Message:  "Service discovery successful",
		Services: convertToServiceInfo(selectedService),
	})
}

// ServiceDeregisterHandler 服务注销处理器
func ServiceDeregisterHandler(c *gin.Context) {
	var req status.ServiceDeregisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, status.ServiceDeregisterResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	// 从 ServiceID 解析服务类型，这里简化处理
	// 实际应用中可能需要从其他途径获取服务类型
	serviceType := "unknown" // 实际实现中需要解析
	serviceID := req.ServiceID

	if err := registryInstance.DeregisterService(serviceType, serviceID); err != nil {
		c.JSON(http.StatusInternalServerError, status.ServiceDeregisterResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, status.ServiceDeregisterResponse{
		Code:    200,
		Message: "Service deregistered successfully",
	})
}

// ServiceHeartbeatHandler 服务心跳处理器
func ServiceHeartbeatHandler(c *gin.Context) {
	var req status_def.ServiceHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, status_def.ServiceHeartbeatResponse{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	// 心跳信息会通过租约机制自动更新，这里只是接收心跳
	c.JSON(http.StatusOK, status_def.ServiceHeartbeatResponse{
		Code:    200,
		Message: "Heartbeat received",
	})
}

// HealthCheckHandler 健康检查处理器
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"code":   200,
	})
}

// CleanupExpiredServices 清理过期服务
func CleanupExpiredServices() {
	// etcd 本身会自动清理过期的键，这里可以做额外的清理工作
}

// 辅助函数：转换内部 ServiceInfo 为 proto ServiceInfo
func convertToServiceInfo(internal *status_model.ServiceInfo) *status_def.ServiceInfo {
	return &status_def.ServiceInfo{
		ServiceName: internal.ServiceName,
		ServiceType: internal.ServiceType,
		ServiceId:   internal.ServiceId,
		Protocol:    internal.Protocol,
		GRPCAddress: internal.GRPCAddress,

		Env:            "default",
		Tags:           []string{internal.IdC},
		HealthCheckUrl: "",
		Weight:         100,
		Status:         "healthy",
		LastHeartbeat:  internal.LastHeartbeat,
	}
}

// ServiceDiscoveryByNameHandler 服务发现处理器（按名称）
func ServiceDiscoveryByNameHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Not implemented",
	})
}

// ServiceDiscoveryByMetadataHandler 服务发现处理器（按元数据）
func ServiceDiscoveryByMetadataHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Not implemented",
	})
}

// ServiceStatusQueryHandler 服务状态查询处理器
func ServiceStatusQueryHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Not implemented",
	})
}

// ServiceStatusOverviewHandler 服务状态概览处理器
func ServiceStatusOverviewHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Not implemented",
	})
}
