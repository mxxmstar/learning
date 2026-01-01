package model

import "time"

type ServiceInfo struct {
	Id            string `json:"id" redis:"id"`                         // 服务唯一标识
	Type          string `json:"type" redis:"type"`                     // 服务类型
	Host          string `json:"host" redis:"host"`                     // 服务主机地址
	Port          int    `json:"port" redis:"port"`                     // 服务端口号
	Metadata      string `json:"metadata" redis:"metadata"`             // 额外元数据
	TTLSeconds    int64  `json:"ttl_seconds" redis:"ttl_seconds"`       // 心跳超时时间
	LastHeartbeat int64  `json:"last_heartbeat" redis:"last_heartbeat"` // 最后心跳时间
	IdC           string `json:"idc" redis:"idc"`                       // 数据中心Id
}

// metadata 中存储的内容：
// MaxLoad       int    `json:"max_load"`       // 最大负载
// Load          int    `json:"load"`           // 当前负载
// UpdatedAt     int64  `json:"updated_at"`     // 最后更新时间
// Version       string `json:"version"`        // 服务版本

// GetRedisKey 返回服务在 Redis 中存储的键
func (s *ServiceInfo) GetRedisKey() string { return "service:" + s.Type + ":" + s.Id }

// IsExpired 检查服务是否已过期
func (s *ServiceInfo) IsExpired() bool { return s.LastHeartbeat+s.TTLSeconds < GetCurrentTimestamp() }

func GetCurrentTimestamp() int64 { return time.Now().Unix() }
