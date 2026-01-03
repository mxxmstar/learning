package status_config

import (
	"fmt"
	"log"

	"github.com/mxxmstar/learning/pkg/config"
)

type StatusServiceConfig struct {
	CleanupInterval     int `mapstructure:"cleanup_interval"`      // 清理间隔(秒) （过期服务，维护服务列表）
	HeartbeatTimeout    int `mapstructure:"heartbeat_timeout"`     // 心跳超时(秒)
	HealthCheckInterval int `mapstructure:"health_check_interval"` // 健康检查间隔(秒)
	ServiceExpireTime   int `mapstructure:"service_expire_time"`   // 服务过期时间(秒)
}

type Config struct {
	ServerConfig  *config.ServerConfig  `mapstructure:"server"` // 服务器配置
	StatusService StatusServiceConfig   `mapstructure:"status_service"`
	Database      config.DatabaseConfig `mapstructure:"database"`
	Redis         config.RedisConfig    `mapstructure:"redis"`
}

func Init() (*Config, error) {
	baseCfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerConfig: &baseCfg.Server,
		StatusService: StatusServiceConfig{
			CleanupInterval:     300,
			HeartbeatTimeout:    60,
			HealthCheckInterval: 30,
			ServiceExpireTime:   120,
		},
		Database: baseCfg.Database,
		Redis:    baseCfg.Redis,
	}

	err = config.Reload(cfg)
	if err != nil {
		return nil, err
	}
	// 打印配置
	log.Printf("Config: %+v\n", cfg)
	return cfg, nil
}

// GetStatusServer 获取status server自己的地址信息
func (c *Config) GetStatusServerHttpAddr() string {
	addr := fmt.Sprintf("%s:%d", c.ServerConfig.StatusServer.HttpConfig.Host, c.ServerConfig.StatusServer.HttpConfig.Port)
	return addr
}

func (c *Config) GetStatusServerGRPCAddr() string {
	addr := fmt.Sprintf("%s:%d", c.ServerConfig.StatusServer.GRPCConfig.Host, c.ServerConfig.StatusServer.GRPCConfig.Port)
	return addr
}
