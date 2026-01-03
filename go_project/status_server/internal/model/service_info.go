package model

import "time"

type ServiceInfo struct {
	ServiceName    string            `json:"service_name"`                          // 服务名称，如 gate_server_1
	ServiceId      string            `json:"service_id"`                            // 服务Id
	Protocol       []string          `json:"protocol"`                              // 服务协议，如：grpc, http
	GRPCAddress    *GRPCAddress      `json:"grpc_address"`                          // grpc服务地址
	HTTPAddress    *HTTPAddress      `json:"http_address"`                          // http服务地址
	Env            string            `json:"env"`                                   // 环境，如：prod, test
	Metadata       map[string]string `json:"metadata,omitempty"`                    // 元数据，存储详细的服务信息
	HealthCheckUrl string            `json:"health_check_url"`                      // 健康检查地址
	Weight         int               `json:"weight"`                                // 权重
	Status         string            `json:"status"`                                // 服务状态
	TTLSeconds     int64             `json:"ttl_seconds" redis:"ttl_seconds"`       // 心跳超时允许时间
	LastHeartbeat  int64             `json:"last_heartbeat" redis:"last_heartbeat"` // 最后心跳时间
	IdC            string            `json:"idc" redis:"idc"`                       // 数据中心Id
}

type GRPCAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HTTPAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// metadata 中存储的内容：
// MaxLoad       int    `json:"max_load"`       // 最大负载
// Load          int    `json:"load"`           // 当前负载
// UpdatedAt     int64  `json:"updated_at"`     // 最后更新时间
// Version       string `json:"version"`        // 服务版本

// IsExpired 检查服务是否已过期
func (s *ServiceInfo) IsExpired() bool { return s.LastHeartbeat+s.TTLSeconds < GetCurrentTimestamp() }

func GetCurrentTimestamp() int64 { return time.Now().Unix() }
