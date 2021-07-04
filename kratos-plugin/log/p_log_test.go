package log

import (
	"context"
	"errors"
	zap_log "github.com/laxiaohong/lollipop/zap-log"
	"github.com/laxiaohong/lollipop/zap-log/init/config"
	"go.uber.org/zap"
	"testing"
	"time"
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

func TestNewCoreLogger_Example(_ *testing.T) {
	var outputConsole = new(bool)
	*outputConsole = true
	zapLogger := zap_log.NewZapLogger(&config.ZapLoggerConfig{Path: "logs", Level: "debug", DebugModeOutputConsole: outputConsole})
	defer zapLogger.Sync()
	var hooks []Option
	hooks = append(hooks, func(ctx context.Context) (field zap.Field, ok bool) {
		traceId := getTraceId(ctx)
		return zap.String("trace_id", traceId), traceId != ""
	})
	hooks = append(hooks, func(ctx context.Context) (field zap.Field, ok bool) {
		return zap.String("ts", time.Now().Format("20060102-15:04:05:06.999")), true
	})
	logger := NewCoreLogger(zapLogger, hooks...)
	var err error = errors.New("div by zero")

	logger.WithContext(context.Background()).Errorf("data error %+v", err)
}

func TestMain(t *testing.M) {
	t.Run()

}
