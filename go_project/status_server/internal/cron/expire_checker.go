package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/mxxmstar/learning/pkg/store/redis"
	"github.com/mxxmstar/learning/status_server/status_config"
)

var (
	Config *status_config.Config // 全局变量 存储配置实例
	store  *redis.RedisClient    // 存储 Redis 实例
)

func StartExpireChecker(config *status_config.Config, redisClient *redis.RedisClient) {
	Config = config
	store = redisClient

	ticker := time.NewTicker(time.Duration(Config.StatusService.CleanupInterval) * time.Second)

	go func() {
		for range ticker.C {
			checkExpiredServices()
		}
	}()
}

func checkExpiredServices() {
	types := []string{"chat", "video", "verify", "file"}
	ctx := context.Background()
	for _, t := range types {
		listKey := fmt.Sprintf("services:%s", t)

		ids, _ := store.ZRange(ctx, listKey, 0, -1).Result()
		for _, id := range ids {
			key := fmt.Sprintf("service:%s:%s", t, id)
			last, _ := store.HGet(ctx, key, "last_heartbeat").Int64()

			if time.Now().Unix()-last > int64(Config.StatusService.HeartbeatTimeout) {
				store.Del(ctx, key)
				store.ZRem(ctx, listKey, id)
				fmt.Println("Removed expired service:", id)
			}
		}
	}
}
