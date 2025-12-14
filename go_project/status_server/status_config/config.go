package status_config

import (
	"log"

	"github.com/mxxmstar/learning/pkg/config"
)

type StatusServiceConfig struct {
	CleanupInterval  int `mapstructure:"cleanup_interval"`
	HeartbeatTimeout int `mapstructure:"heartbeat_timeout"`
}

type StatusServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	StatusServer   StatusServerConfig  `mapstructure:"status_server"`
	StatusService  StatusServiceConfig `mapstructure:"status_service"`
	*config.Config `mapstructure:",squash"`
}

func Init() (*Config, error) {
	baseCfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Config: baseCfg,
	}

	err = config.Reload(cfg)
	if err != nil {
		return nil, err
	}
	// 打印配置
	log.Printf("Config: %+v\n", cfg)
	return cfg, nil
}
