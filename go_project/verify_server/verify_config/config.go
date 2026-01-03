package verify_config

import (
	"log"

	"github.com/mxxmstar/learning/pkg/config"
)

type VerifyServiceConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`     // jwt密钥
	TokenLifeTime int    `mapstructure:"token_lifetime"` // token有效期
	RefreshToken  bool   `mapstructure:"refresh_token"`  // 是否允许刷新token
}

type Config struct {
	ServerConfig  *config.ServerConfig  `mapstructure:"server"`         // 	服务器配置
	Database      config.DatabaseConfig `mapstructure:"database"`       //数据库配置
	Redis         config.RedisConfig    `mapstructure:"redis"`          // redis配置
	VerifyService VerifyServiceConfig   `mapstructure:"verify_service"` // 验证服务特定的配置
}

func Init() (*Config, error) {
	baseCfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerConfig: &baseCfg.Server,
		Database:     baseCfg.Database,
		Redis:        baseCfg.Redis,
		// TODO: 从配置文件中读取,密钥由密钥管理服务生成
		VerifyService: VerifyServiceConfig{
			JWTSecret:     "secret",
			TokenLifeTime: 86400,
			RefreshToken:  true,
		},
	}

	err = config.Reload(cfg)
	if err != nil {
		return nil, err
	}
	// 打印配置
	log.Printf("Config: %+v\n", cfg)
	return cfg, nil
}
