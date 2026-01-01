package logger

import (
	"context"
	"fmt"
	"time"
)

// ConnectionContext 用于存储连接相关的信息
type ConnectionContext struct {
	ConnId    string
	TraceId   string // 链路追踪Id 16字符
	StartTime time.Time
}

type ConnectionContextKey struct{}

// 从上下文 context 中获取 ConnectionContext
func GetConnectionContext(ctx context.Context) *ConnectionContext {
	if connCtx, ok := ctx.Value(ConnectionContextKey{}).(*ConnectionContext); ok {
		return connCtx
	}
	return nil
}

// 将 ConnectionContext 存储到上下文 context 中
func WithConnectionContext(ctx context.Context, connId string) context.Context {
	connCtx := &ConnectionContext{
		ConnId:    connId,
		TraceId:   NewTraceId(),
		StartTime: time.Now(),
	}
	return context.WithValue(ctx, ConnectionContextKey{}, connCtx)
}

// 格式化日志上下文信息
func FormatLogContext(ctx context.Context) string {
	connCtx := GetConnectionContext(ctx)
	if connCtx == nil {
		return ""
	}
	return fmt.Sprintf("[conn:%s][traceId:%s]", connCtx.ConnId, connCtx.TraceId)
}
