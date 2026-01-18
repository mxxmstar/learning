package status_def

import (
	status_model "github.com/mxxmstar/learning/pkg/model/status"
)

/**
 * @Description: 定义其他服务向 status_server 发送的请求和响应
 * 	包含：服务注册、服务注销、服务心跳、服务发现、服务状态查询、服务健康检查
 **/

// ServiceRegisterRequest 服务注册请求
type ServiceRegisterRequest struct {
	ServiceName    string            `json:"service_name"`         // 服务名称，如 gate_server_1
	ServiceType    string            `json:"service_type"`         // 服务类型，如：gate, verify
	ServiceId      string            `json:"service_id,omitempty"` // 服务Id
	Protocol       []string          `json:"protocol"`             // 服务协议，如：["grpc", "http"]
	GRPCAddress    *GRPCAddress      `json:"grpc_address"`         // grpc服务地址
	HTTPAddress    *HTTPAddress      `json:"http_address"`         // http服务地址
	Env            string            `json:"env"`                  // 环境，如：prod, test
	Tags           []string          `json:"tags,omitempty"`       // 标签，如：["region=cn-hangzhou", "zone=hangzhou-b"]
	Metadata       map[string]string `json:"metadata,omitempty"`   // 元数据，存储详细的服务信息
	HealthCheckUrl string            `json:"health_check_url"`     // 健康检查地址
	Weight         int               `json:"weight"`               // 权重
	Enable         bool              `json:"enable"`               // 是否启用
	Idc            string            `json:"idc"`                  // 机房
}

type GRPCAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HTTPAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (req *GRPCAddress) ConvertToStatusServiceInfo() *status_model.GRPCAddress {
	grpc := &status_model.GRPCAddress{
		Host: req.Host,
		Port: req.Port,
	}
	return grpc
}

func (req *HTTPAddress) ConvertToStatusServiceInfo() *status_model.HTTPAddress {
	http := &status_model.HTTPAddress{
		Host: req.Host,
		Port: req.Port,
	}
	return http
}
func (req *ServiceRegisterRequest) ConvertToStatusServiceInfo() *status_model.ServiceInfo {
	status := "online"
	if req.Enable {
		status = "active"
	}
	// 转换为内部 ServiceInfo 结构
	serviceInfo := &status_model.ServiceInfo{
		ServiceName:    req.ServiceName,
		ServiceType:    req.ServiceType,
		ServiceId:      req.ServiceId,
		Protocol:       req.Protocol,
		GRPCAddress:    req.GRPCAddress.ConvertToStatusServiceInfo(),
		HTTPAddress:    req.HTTPAddress.ConvertToStatusServiceInfo(),
		Env:            req.Env,
		Metadata:       req.Metadata,
		HealthCheckUrl: req.HealthCheckUrl,
		Weight:         req.Weight,
		Status:         status,
		TTLSeconds:     status_model.HeartbeatInterval,
		IdC:            req.Idc,
	}
	return serviceInfo
}

// ServiceRegisterResponse 服务注册响应
type ServiceRegisterResponse struct {
	ServiceId string `json:"service_id"` // 注册成功的服务Id
	Code      int    `json:"code"`       // 状态码
	Message   string `json:"message"`    // 状态信息
}

// ServiceDeregisterRequest 服务注销请求
type ServiceDeregisterRequest struct {
	ServiceId string `json:"service_id"` // 服务Id
}

type ServiceDeregisterResponse struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 状态信息
}

// ServiceHeartbeatRequest 服务心跳请求
type ServiceHeartbeatRequest struct {
	ServiceId string `json:"service_id"` // 服务Id
	Status    string `json:"status"`     // 服务状态
	Timestamp int64  `json:"timestamp"`  // 服务心跳时间
}

type ServiceHeartbeatResponse struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 状态信息
}

// ServiceDiscoveryRequest 服务发现请求
type ServiceDiscoveryRequest struct {
	ServiceName string            `json:"service_name"`       // 服务名称
	Tags        []string          `json:"tags,omitempty"`     // 过滤标签
	Metadata    map[string]string `json:"metadata,omitempty"` // 过滤元数据
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceName    string            `json:"service_name"`       // 服务名称，如 gate_server_1
	ServiceType    string            `json:"service_type"`       // 服务类型，如：gate, verify
	ServiceId      string            `json:"service_id"`         // 服务Id
	Protocol       []string          `json:"protocol"`           // 服务协议，如：grpc, http
	GRPCAddress    *GRPCAddress      `json:"grpc_address"`       // grpc服务地址
	HTTPAddress    *HTTPAddress      `json:"http_address"`       // http服务地址
	Env            string            `json:"env"`                // 环境，如：prod, test
	Tags           []string          `json:"tags,omitempty"`     // 标签，如：["region=cn-hangzhou", "zone=hangzhou-b"]
	Metadata       map[string]string `json:"metadata,omitempty"` // 元数据，存储详细的服务信息
	HealthCheckUrl string            `json:"health_check_url"`   // 健康检查地址
	Weight         int               `json:"weight"`             // 权重
	Status         string            `json:"status"`             // 服务状态
	LastHeartbeat  int64             `json:"last_heartbeat"`     // 最后一次心跳时间
}

// ServiceDiscoveryResponse 服务发现响应
type ServiceDiscoveryResponse struct {
	Code     int            `json:"code"`     // 状态码
	Message  string         `json:"message"`  // 状态信息
	Services []*ServiceInfo `json:"services"` // 服务列表
}

// ServiceStatusRequest 服务状态查询请求
type ServiceStatusRequest struct {
	ServiceName string `json:"service_name,omitempty"` // 服务名，为空则查询所有服务
	ServiceId   string `json:"service_id,omitempty"`   // 服务Id，为空则查询所有服务
}

// ServiceStatusResponse 服务状态查询响应
type ServiceStatusResponse struct {
	Code     int            `json:"code"`     // 状态码
	Message  string         `json:"message"`  // 状态信息
	Services []*ServiceInfo `json:"services"` // 服务列表
}

// HealthCheckRequest 服务健康检查请求
type HealthCheckRequest struct {
	ServiceId string `json:"service_id"` // 服务Id
}
type HealthCheckResponse struct {
	Code      int    `json:"code"`       // 状态码
	Message   string `json:"message"`    // 状态信息
	Status    string `json:"status"`     // 服务状态
	ServiceId string `json:"service_id"` // 服务Id
}

type ServiceDiscoveryByTagsRequest struct {
	ServiceName string            `json:"service_name"`       // 服务名
	Tags        []string          `json:"tags,omitempty"`     // 标签
	Metadata    map[string]string `json:"metadata,omitempty"` // 元数据
	Strategy    string            `json:"strategy"`           // 负载均衡策略
}

type ServiceDiscoveryByTagsResponse struct {
	Code     int          `json:"code"`    // 状态码
	Message  string       `json:"message"` // 状态信息
	Services *ServiceInfo `json:"service"` // 服务实例
}
