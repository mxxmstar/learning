package config

import "time"

// 服务器配置结构
type ProjectConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`
}

// 数据库实例配置结构
type DatabaseInstanceConfig struct {
	DSN          string        `mapstructure:"dsn"`
	CharSet      string        `mapstructure:"charset" default:"utf8mb4"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout" default:"5s"`  // 连接超时时间
	ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"3s"`  // 读取超时时间
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"3s"` // 写入超时时间
}

// 数据库连接池结构
type DatabasePoolConfig struct {
	// 最大连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// 连接最大生命周期 开始创建-最后使用时间 超过该时间强制关闭连接
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	// 连接最大空闲时间 超过该时间空闲连接将被关闭，避免长时间占用数据库连接资源
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

type RedisInstanceConfig struct {
	DSN         string        `mapstructure:"dsn"`
	Addr        string        `mapstructure:"addr"`
	Password    string        `mapstructure:"password"`
	DB          int           `mapstructure:"db"`
	DialTimeout time.Duration `mapstructure:"dial_timeout" default:"5s"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"` // Kafka broker 地址列表
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type MessageQueueConfig struct {
	Type  string              `mapstructure:"type"`  // 消息队列类型，如 "kafka" "rabbitmq" "redis"
	Redis RedisInstanceConfig `mapstructure:"redis"` // redis pub/sub 配置
	Kafka KafkaConfig         `mapstructure:"kafka"` // kafka 配置
}

type RedisPoolConfig struct {
	PoolSize int `mapstructure:"pool_size" default:"10"`
	// 最小空闲连接数
	MinIdleConns int `mapstructure:"min_idle_conns" default:"5"`
	// 命令执行失败时的最大重试次数
	MaxRetries int `mapstructure:"max_retries" default:"3"`
}

type RedisConfig struct {
	Mode string `mapstructure:"mode" default:"cluster"` // redis 模式，如 "cluster" "standalone"
	// 单实例
	Standalone RedisInstanceConfig `mapstructure:"standalone"`
	// 集群实例
	Cluster []RedisInstanceConfig `mapstructure:"cluster"`
	// 其它实例
	Session RedisInstanceConfig `mapstructure:"session"` // 会话存储
	Cache   RedisInstanceConfig `mapstructure:"cache"`   // 缓存存储
}

// 主配置结构
type Config struct {
	// 环境变量
	ProjectConfig ProjectConfig `mapstructure:"PROJECT"`

	Database     DatabaseConfig     `mapstructure:"DATABASE"`
	Redis        RedisConfig        `mapstructure:"REDIS"`
	MessageQueue MessageQueueConfig `mapstructure:"MESSAGE_QUEUE"`
	// 添加其它需要的配置项
}
