package config

import "time"

type RedisInstanceConfig struct {
	DSN         string        `mapstructure:"dsn"`
	Addr        string        `mapstructure:"addr"`
	Password    string        `mapstructure:"password"`
	DB          int           `mapstructure:"db"`
	DialTimeout time.Duration `mapstructure:"dial_timeout" default:"5s"`
}

type RedisPoolConfig struct {
	PoolSize int `mapstructure:"pool_size" default:"10"`
	// 最小空闲连接数
	MinIdleConns int `mapstructure:"min_idle_conns" default:"5"`
	// 命令执行失败时的最大重试次数
	MaxRetries int `mapstructure:"max_retries" default:"3"`
}

type RedisConfig struct {
	Mode       string                `mapstructure:"mode" default:"standalone"` // redis 模式，如 "cluster" "standalone"
	Standalone RedisInstanceConfig   `mapstructure:"standalone"`                // 单实例
	Cluster    []RedisInstanceConfig `mapstructure:"cluster"`                   // 集群实例
	Usage      RedisUsageConfig      `mapstructure:"usage"`                     // 使用方式
	Pool       RedisPoolConfig       `mapstructure:"pool"`                      // 连接池配置
}

type RedisUsageConfig struct {
	Session RedisInstanceConfig `mapstructure:"session"` // 会话存储
	Cache   RedisInstanceConfig `mapstructure:"cache"`   // 缓存存储
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
