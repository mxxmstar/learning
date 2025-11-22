package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init 初始化配置
func Init() (*Config, error) {
	// 环境变量
	viper.SetDefault("server.name", "learn_mx")
	viper.SetDefault("server.env", "dev")

	// 数据库配置
	viper.SetDefault("database.type", "mysql")
	viper.SetDefault("database.auth_db.dsn", "root:123456@tcp(127.0.0.1:3306)/learn_mysql?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("database.chat_db.dsn", "root:123456@tcp(127.0.0.1:3306)/learn_mysql?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("database.file_db.dsn", "root:123456@tcp(127.0.0.1:3306)/learn_mysql?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("database.video_db.dsn", "root:123456@tcp(127.0.0.1:3306)/learn_mysql?charset=utf8mb4&parseTime=True&loc=Local")

	viper.SetDefault("database.pool.max_open_conns", 50)
	viper.SetDefault("database.pool.max_idle_conns", 10)
	viper.SetDefault("database.pool.conn_max_lifetime", "300s")
	viper.SetDefault("database.pool.conn_max_idle_time", "120s")

	viper.SetDefault("database.auto_migrate", true)

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
