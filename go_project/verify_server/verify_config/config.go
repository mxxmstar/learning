package verify_config

import (
	"log"

	"github.com/mxxmstar/learning/pkg/config"
)

type VerifyServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	// JWTSecret     string `mapstructure:"jwt_secret"`
	// TokenLifeTime int    `mapstructure:"token_lifetime"`
	// LogLevel string `mapstructure:"log_level"`
}

type Config struct {
	VerifyServer   VerifyServerConfig `mapstructure:"verify_server"`
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
