package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotAcquired = errors.New("lock not acquired")
	ErrLockNotReleased = errors.New("lock not released")
	ErrLockNotOwned    = errors.New("lock is not owned by this instance")
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{client: client}
}

// 设置键值对，过期时间为expiration(可选)
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := rc.client.Set(ctx, key, value, expiration).Result()
	return err
}

// 获取键对应的值
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *RedisClient) Del(ctx context.Context, keys ...string) error {
	return rc.client.Del(ctx, keys...).Err()
}

func (rc *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rc.client.Expire(ctx, key, expiration).Err()
}

// Eval 执行 Lua 脚本
func (rc *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return rc.client.Eval(ctx, script, keys, args...)
}

// ---------- 分布式锁 ----------
type DistributedLock struct {
	client   *RedisClient
	key      string
	value    string // 用于标识锁持有者
	duration time.Duration
}

func NewDistributedLock(client *RedisClient, key string, duration time.Duration) (*DistributedLock, error) {
	value, err := generateLockValue()
	if err != nil {
		return nil, err
	}
	return &DistributedLock{
		client:   client,
		key:      key,
		value:    value,
		duration: duration,
	}, nil
}

// 尝试获取锁 非阻塞
func (l *DistributedLock) TryLock(ctx context.Context) error {
	ok, err := l.client.client.SetNX(ctx, l.key, l.value, l.duration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrLockNotAcquired
	}
	return nil
}

func (l *DistributedLock) Lock(ctx context.Context, retryInterval time.Duration) error {
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := l.TryLock(ctx); err == nil {
				return nil
			} else if !errors.Is(err, ErrLockNotAcquired) {
				return err
			}
			// 锁未获取成功，继续尝试
		}
	}
}

// Unlock 释放锁（安全释放：只有持有者才能释放）
func (l *DistributedLock) Unlock(ctx context.Context) error {
	script := `
				if redis.call("get", KEYS[1]) == ARGV[1] then
					return redis.call("del", KEYS[1])
				else
					return 0
				end
			`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return ErrLockNotOwned
	}
	return nil
}

// 自动续期（可选，配合 goroutine 使用）
func (l *DistributedLock) StartAutoRefresh(ctx context.Context, refreshInterval time.Duration) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 尝试延长锁时间（仅当仍持有锁时）
				_, err := l.client.Eval(ctx, `
											if redis.call("get", KEYS[1]) == ARGV[1] then
												return redis.call("pexpire", KEYS[1], ARGV[2])
											else
												return 0
											end
										`, []string{l.key}, l.value, int64(l.duration/time.Millisecond)).Result()
				if err != nil {
					// 可记录日志，但不中断
					continue
				}
			}
		}
	}()
	return done
}

func generateLockValue() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
