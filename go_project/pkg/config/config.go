package config

import "time"

// 服务器配置结构
type ServerConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

// 数据库实例配置结构
type DatabaseInstanceConfig struct {
	DSN string `mapstructure:"dsn"`
}

// 数据库连接池结构
type DatabasePoolConfig struct {
	// 最大连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 连接最大生命周期
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	// 连接最大空闲时间
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// 数据库配置结构
type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	// 认证库 所有服务共享
	AuthDB DatabaseInstanceConfig `mapstructure:"auth_db"`
	// 聊天库（ChatServer）
	ChatDB DatabaseInstanceConfig `mapstructure:"chat_db"`
	// 文件服务库（FileServer）
	FileDB DatabaseInstanceConfig `mapstructure:"file_db"`
	// 视频/RTC库（VideoServer + RTCServer）
	VideoDB DatabaseInstanceConfig `mapstructure:"video_db"`
	// 数据库连接池配置
	Pool DatabasePoolConfig `mapstructure:"pool"`
	// 是否自动迁移数据库表结构
	AutoMigrate bool `mapstructure:"auto_migrate"`
}

type RedisConfig struct {
}

// 主配置结构
type Config struct {
	// 环境变量
	Env string `mapstructure:"ENV"`

	Database DatabaseConfig `mapstructure:"DATABASE"`

	// 添加其它需要的配置项
}
