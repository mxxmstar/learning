package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init() (*Config, error) {
	viper.SetDefault("ENV", "dev") // 设置默认环境为 dev

	// 支持从 YAML 文件加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")        // 配置文件所在目录
	viper.AddConfigPath("./config") // 配置文件所在目录
	// viper.AddConfigPath("/etc/app") // 可选：系统级配置目录

	// 自动从环境变量中读取配置
	viper.AutomaticEnv()

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %s", err)
	}

	// 绑定配置到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("fatal error unmarshal config file: %s", err)
	}

	return &cfg, nil
}
