package logger

import (
	"context"
	"testing"
	"time"

	// "github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func setupTestLogger(t *testing.T) (func(), *observer.ObservedLogs) {
	// 创建观察者核心，用于捕获日志输出
	observerCore, observedLogs := observer.New(zap.InfoLevel)

	// 创建带有观察者核心的测试日志器
	testLogger := zap.New(observerCore)

	// 保存原始全局日志器
	originalLogger := zap.L()

	// 替换为测试日志器
	zap.ReplaceGlobals(testLogger)

	// 返回清理函数和观察者日志
	return func() {
		zap.ReplaceGlobals(originalLogger) // 恢复原始日志器
	}, observedLogs
}

func TestInitLogger(t *testing.T) {
	// 保存原始全局日志器
	originalLogger := zap.L()
	defer func() {
		zap.ReplaceGlobals(originalLogger) // 恢复原始日志器
	}()

	InitLogger() // 初始化日志器

	// 验证全局日志器已经被替换
	currentLogger := zap.L()
	assert.NotEqual(t, originalLogger, currentLogger, "全局日志器没有被替换")

	// 验证日志配置正确
	logger := currentLogger.WithOptions(zap.AddCallerSkip(-1))
	logger.Info("测试日志")
}

func TestNewTraceId(t *testing.T) {
	// 测试 NewTraceId 函数
	traceId1 := NewTraceId()
	assert.NotEmpty(t, traceId1, "生成的 traceId 为空")
	traceId2 := NewTraceId()
	assert.NotEmpty(t, traceId2, "生成的 traceId 为空")
	assert.NotEqual(t, traceId1, traceId2, "生成的 traceId 重复")
}

func TestWithTraceId(t *testing.T) {
	// 创建基础上下文
	ctx := context.Background()
	// 添加 traceId 到上下文 ctx 中
	ctxWithTraceId := WithTraceId(ctx)
	// 验证 traceId 已添加到上下文
	traceId := GetTraceId(ctxWithTraceId)

	assert.NotEqual(t, "unknown", traceId, "traceId 为 unknown")

	// 测试没有 TraceId 的情况
	ctxWithoutTraceId := context.Background()
	traceId = GetTraceId(ctxWithoutTraceId)
	assert.Equal(t, "unknown", traceId, "traceId 不为 unknown")
}

func TestFormatLog(t *testing.T) {
	// 设置测试环境
	cleanup, observedLogs := setupTestLogger(t)
	defer cleanup() // 测试结束后清理日志器

	// 创建带 TraceId 的上下文
	ctxWithTraceId := WithTraceId(context.Background())

	// 测试 FormatLog 函数
	FormatLog(ctxWithTraceId, "info", "测试日志", zap.String("key", "value"))

	// {"level":"INFO","ts":"2025-09-25T14:30:00.123+08:00","caller":"logger_test.go:80",
	// "msg":"测试日志","trace_id":"a1b2c3d4e5f6g7h8","context":"[Conn:conn123][Trace:a1b2c3d4e5f6g7h8]",
	// "key":"value"}

	// 验证日志输出
	// logs := observedLogs.FilterMessage("测试日志").All()
	logs := observedLogs.TakeAll()
	assert.Len(t, logs, 1, "日志中没有包含 测试日志 消息")

	// 验证日志字段
	log := logs[0]
	assert.Equal(t, "value", log.ContextMap()["key"], "key 字段错误")
	assert.NotEmpty(t, log.ContextMap()["traceId"], "traceId 字段为空")
}

func TestLogAuth(t *testing.T) {
	// 设置测试环境
	cleanup, observedLogs := setupTestLogger(t)
	defer cleanup() // 测试结束后清理日志器

	// 创建带 TraceId 的上下文
	ctxWithTraceId := WithTraceId(context.Background())

	// 测试 LogAuth 函数
	LogAuth(ctxWithTraceId, "login", true, "用户登录成功")

	// 验证日志输出
	logs := observedLogs.FilterMessage("auth").All()
	assert.Len(t, logs, 1, "日志中没有包含 auth 消息")

	found := false
	// 验证日志字段
	for _, log := range logs {
		if log.Message == "auth" {
			found = true
			assert.Equal(t, "login", log.ContextMap()["action"], "action 字段错误")
			assert.Equal(t, true, log.ContextMap()["success"], "success 字段错误")
			assert.Equal(t, "用户登录成功", log.ContextMap()["msg"], "消息内容错误")
			assert.NotEmpty(t, log.ContextMap()["traceId"], "traceId 字段为空")
			break
		}
	}
	assert.True(t, found, "没有找到 auth 消息")
}

func TestLogRouter(t *testing.T) {
	// 设置测试环境
	cleanup, observedLogs := setupTestLogger(t)
	defer cleanup() // 测试结束后清理日志器

	// 创建带 TraceId 的上下文
	ctxWithTraceId := WithTraceId(context.Background())

	// 测试 LogRouter 函数
	duration := 100 * time.Millisecond
	LogRouter(ctxWithTraceId, "/api/v1/hello", duration)

	// 验证日志输出
	logs := observedLogs.FilterMessage("router").All()
	assert.Len(t, logs, 1, "日志中没有包含 router 消息")

	// 验证日志字段
	log := logs[0]
	assert.Equal(t, "/api/v1/hello", log.ContextMap()["route"], "route 字段错误")
	assert.Equal(t, duration, log.ContextMap()["duration"], "duration 字段错误")
	assert.NotEmpty(t, log.ContextMap()["traceId"], "traceId 字段为空")
}
