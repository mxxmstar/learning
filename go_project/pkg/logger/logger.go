package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化日志记录器
func InitLogger() {
	// 生产环境日志配置
	config := zap.NewProductionConfig()
	// 时间格式
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志级别大写
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 调用者信息的键名
	config.EncoderConfig.CallerKey = "caller"
	// 函数信息的键名
	config.EncoderConfig.FunctionKey = "function"

	// 创建日志器
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	// 替换全局日志器
	zap.ReplaceGlobals(logger)
}

func FormatLog(ctx context.Context, level string, msg string, fields ...zap.Field) {
	// 获取 traceID
	traceID := GetTraceID(ctx)

	// 添加 traceID 到日志字段
	fields = append(fields, zap.String("traceID", traceID))

	// 获取日志上下文信息
	logContext := FormatLogContext(ctx)
	// 添加日志上下文到字段
	if logContext != "" {
		fields = append(fields, zap.String("context", logContext))
	}

	// 记录日志
	switch level {
	case "debug", "DEBUG":
		zap.L().Debug(msg, fields...)
	case "info", "INFO":
		zap.L().Info(msg, fields...)
	case "warn", "WARN":
		zap.L().Warn(msg, fields...)
	case "error", "ERROR":
		zap.L().Error(msg, fields...)
	case "fatal", "FATAL":
		zap.L().Fatal(msg, fields...)
	default:
		zap.L().Info(msg, fields...)
	}
}

// 关键路径日志记录
func LogAuth(ctx context.Context, action string, success bool, msg string) {
	fields := []zap.Field{
		zap.String("action", action),
		zap.Bool("success", success),
	}

	if msg != "" {
		fields = append(fields, zap.String("msg", msg))
	}
	FormatLog(ctx, "info", "auth", fields...)
}

func LogRouter(ctx context.Context, route string, duration time.Duration) {
	fields := []zap.Field{
		zap.String("route", route),
		zap.Duration("duration", duration),
	}
	FormatLog(ctx, "info", "router", fields...)
}
