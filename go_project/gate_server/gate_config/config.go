package gate_config

import (
	"log"

	"github.com/mxxmstar/learning/pkg/config"
)

type VerifyGRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GateServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	// LogLevel string `mapstructure:"log_level"`
}

type Config struct {
	GateServer       GateServerConfig       `mapstructure:"gate_server"`
	VerifyGRPCServer VerifyGRPCServerConfig `mapstructure:"verify_grpc"`
	*config.Config   `mapstructure:",squash"`
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
