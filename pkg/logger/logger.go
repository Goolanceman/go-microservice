package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// Init initializes the global logger only once with level and file path
func Init(level, filePath string) error {
	var initErr error
	once.Do(func() {
		config := zap.NewProductionConfig()

		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
			zapLevel = zapcore.InfoLevel
		}
		config.Level = zap.NewAtomicLevelAt(zapLevel)

		if filePath != "" {
			config.OutputPaths = []string{filePath, "stdout"}
		} else {
			config.OutputPaths = []string{"stdout"}
		}

		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		instance, initErr = config.Build()
	})
	return initErr
}

// FromContext returns a logger with context fields
func FromContext(ctx context.Context) *zap.Logger {
	if instance == nil {
		return zap.NewNop()
	}
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return instance.With(zap.String("request_id", reqID))
	}
	return instance
}

// Info logs an info level message
func Info(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Info(msg, fields...)
	}
}

// Error logs an error level message
func Error(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Error(msg, fields...)
	}
}

// Debug logs a debug level message
func Debug(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Debug(msg, fields...)
	}
}

// Warn logs a warn level message
func Warn(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Warn(msg, fields...)
	}
}

// Fatal logs a fatal level message
func Fatal(msg string, fields ...zap.Field) {
	if instance != nil {
		instance.Fatal(msg, fields...)
	}
}

// Sync flushes the logger
func Sync() error {
	if instance != nil {
		return instance.Sync()
	}
	return nil
}
