package log

import (
	"context"
	zap_log "github.com/laxiaohong/lollipop/zap-log"
	"github.com/laxiaohong/lollipop/zap-log/init/config"
	"go.uber.org/zap"
	"testing"
)

func TestNewCoreLogger(t *testing.T) {
	logger := zap_log.NewZapLogger(&config.ZapLoggerConfig{Path: "logs", Level: "error"})
	defer logger.Sync()
	v := NewCoreLogger(logger)
	_ = v
	v.WithContext(context.TODO()).Info("hello world")
	v.WithContext(context.TODO(), /*WithTraceId(getTraceId(context.TODO()))*/).Infow("hello world", "world", "你好", "世界")
	v.WithContext(context.TODO()).Info("hello", "world", "tao", "lu")
	v.WithContext(context.TODO(), func(ctx context.Context) func(c *loggerPointConfig) {
		return func(c *loggerPointConfig) {
			c.field = append(c.field, zap.String("hello", "world"))
		}
	}).Info("woo")
}

func BenchmarkNewCoreLogger(b *testing.B) {
	zapLogger := zap_log.NewZapLogger(&config.ZapLoggerConfig{Path: "logs", Level: zap_log.LevelInfo})
	defer zapLogger.Sync()
	logger := NewCoreLogger(zapLogger)
	for i := 0; i < b.N; i++ {
		logger.WithContext(context.TODO()).Debug("debug")
		logger.WithContext(context.TODO()).Info("info")
		logger.WithContext(context.TODO()).Warn("warn")
		logger.WithContext(context.TODO()).Error("error")
		zapLogger.Info("hello")
	}
}
