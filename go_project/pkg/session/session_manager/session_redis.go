package session_manager

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mxxmstar/learning/pkg/logger"
	"github.com/mxxmstar/learning/pkg/session"
	"github.com/mxxmstar/learning/pkg/store/redis"
)

var (
	ErrSessionNotFound   = errors.New("session not found in redis")
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionDelFailed  = errors.New("session deleted failed")
	ErrSessionStrMarshal = errors.New("session string marshal failed")
)

type RedisSessionManager struct {
	redisClient *redis.RedisClient
	prefix      string
	defaultTTL  time.Duration
}

func NewRedisSessionManager(redisClient *redis.RedisClient, prefix string, defaultTTL time.Duration) *RedisSessionManager {
	return &RedisSessionManager{
		redisClient: redisClient,
		prefix:      prefix,
		defaultTTL:  defaultTTL,
	}
}

// GetSession 从redis中获取会话
func (m *RedisSessionManager) GetSession(sessionId string) (*session.BaseSession, error) {
	ctx := context.Background()
	key := m.prefix + sessionId

	sessionStr, err := m.redisClient.Get(ctx, key)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	var session session.BaseSession
	err = json.Unmarshal([]byte(sessionStr), &session)
	if err != nil {
		return nil, ErrSessionStrMarshal
	}

	// 检查会话是否过期
	if session.IsExpired() {
		// 会话过期，删除redis中的会话
		err = m.redisClient.Del(ctx, key)
		if err != nil {
			logger.LogAuth(ctx, "redis", false, "delete session failed: "+err.Error())
			return nil, ErrSessionDelFailed
		}
		return nil, ErrSessionExpired
	}

	logger.LogAuth(ctx, "redis", true, "get session: "+sessionId)
	return &session, nil
}

func (m *RedisSessionManager) AddSession(session *session.BaseSession) error {
	ctx := context.Background()
	key := m.prefix + string(session.Id)

	sessionStr, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// 设置redis中的会话
	return m.redisClient.Set(ctx, key, string(sessionStr), m.defaultTTL)
}

func (m *RedisSessionManager) RemoveSession(sessionId string) error {
	ctx := context.Background()
	key := m.prefix + sessionId

	// 删除redis中的会话
	return m.redisClient.Del(ctx, key)
}

// ExtendExpiration 延长会话的过期时间
func (m *RedisSessionManager) ExtendExpiration(sessionId string, ttl time.Duration) error {
	ctx := context.Background()
	key := m.prefix + sessionId

	// 延长redis中的会话的过期时间
	return m.redisClient.Expire(ctx, key, ttl)
}

// GetSessionsByType 获取指定类型的会话（该操作需要遍历所有会话，性能可能受到影响）
func (m *RedisSessionManager) GetSessionsByType(sessionType session.SessionType) ([]*session.BaseSession, error) {
	// 注意：由于Redis的特性，此方法在大量Session情况下性能较差
	// 建议在实际项目中避免频繁调用或者添加额外的索引机制
	// TODO: 实现根据会话类型获取会话的逻辑
	return nil, nil
}

func (m *RedisSessionManager) GetSessionCount() (int, error) {
	// 注意：由于Redis的特性，此方法在大量Session情况下性能较差
	// 建议在实际项目中避免频繁调用或者添加额外的索引机制
	// TODO: 实现获取会话数量的逻辑
	return 0, nil
}
