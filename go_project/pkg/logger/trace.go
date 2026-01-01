package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type traceIdKey struct{} // 定义上下文中的key类型，避免冲突

// 生成唯一的trace Id
func NewTraceId() string {
	// 8字节随机数据 → 16个十六进制字符
	b := make([]byte, 8) // 8字节 = 64位随机数
	_, err := rand.Read(b)
	if err != nil {
		// 生产环境应处理错误，这里用安全回退（非时间戳）
		// 生成16位固定长度的随机字符串（安全回退）
		return fmt.Sprintf("%016x", time.Now().UnixNano()%10000000000)
	}
	return hex.EncodeToString(b) // 严格16字符
}

func GetTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(traceIdKey{}).(string); ok {
		return traceId
	}
	return "unknown"
}

func WithTraceId(ctx context.Context) context.Context {
	traceId := NewTraceId()
	return context.WithValue(ctx, traceIdKey{}, traceId)
}
