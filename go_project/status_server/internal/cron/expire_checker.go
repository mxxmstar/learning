package cron

import (
	"fmt"
	"time"
)

func StartExpireChecker() {
	ticker := time.NewTicker(time.Duration(config.Cfg.Service.CleanupInterval) * time.Second)

	go func() {
		for range ticker.C {
			checkExpiredServices()
		}
	}()
}

func checkExpiredServices() {
	types := []string{"chat", "video", "verify", "file"}

	for _, t := range types {
		listKey := fmt.Sprintf("services:%s", t)

		ids, _ := store.Rdb.ZRange(store.Ctx, listKey, 0, -1).Result()
		for _, id := range ids {
			key := fmt.Sprintf("service:%s:%s", t, id)
			last, _ := store.Rdb.HGet(store.Ctx, key, "last_heartbeat").Int64()

			if time.Now().Unix()-last > int64(config.Cfg.Service.HeartbeatTimeout) {
				store.Rdb.Del(store.Ctx, key)
				store.Rdb.ZRem(store.Ctx, listKey, id)
				fmt.Println("Removed expired service:", id)
			}
		}
	}
}
