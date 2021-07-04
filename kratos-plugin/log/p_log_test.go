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
	v := NewCoreLogger(logger, func(ctx context.Context) (field zap.Field, ok bool) {
		traceId := getTraceId(ctx)
		return zap.String("my_trace_id", traceId), true
	})
	_ = v
	v.WithContext(context.TODO()).Info("hello world")
	v.WithContext(context.TODO() /*WithTraceId(getTraceId(context.TODO()))*/).Infow("hello world", "world", "你好", "世界")
	v.WithContext(context.TODO()).Info("hello", "world", "tao", "lu")
	v.WithContext(context.TODO()).Info()

}

func BenchmarkNewCoreLogger(b *testing.B) {
	zapLogger := zap_log.NewZapLogger(&config.ZapLoggerConfig{Path: "logs", Level: zap_log.LevelDebug})
	defer zapLogger.Sync()
	logger := NewCoreLogger(zapLogger)
	for i := 0; i < b.N; i++ {
		logger.WithContext(context.TODO()).Debug("debug")
		logger.WithContext(context.TODO()).Info("info")
		logger.WithContext(context.TODO()).Warn("warn")
		logger.WithContext(context.TODO()).Error("error")
		//zapLogger.Info("hello")
	}
}

func TestMain(t *testing.M) {
	t.Run()

}
