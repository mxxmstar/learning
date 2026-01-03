package verify_config

import "github.com/mxxmstar/learning/pkg/store/redis"

func InitRedis(cfg *Config) (*redis.RedisClient, error) {
	// 使用配置中的Redis实例来初始化 Redis 客户端
	redisCfg := cfg.Redis.Standalone // redis 单机实例
	client := redis.NewRedisClient(redisCfg.Addr, redisCfg.Password, redisCfg.DB)
	return client, nil
}
