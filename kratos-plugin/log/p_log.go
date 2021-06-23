package log

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"sync"
)

func NewCoreLogger(logger *zap.Logger) *CoreLogger {
	return &CoreLogger{
		logger: logger,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

type CoreLogger struct {
	logger *zap.Logger
	pool   *sync.Pool
}

// Option type
type Option = func(ctx context.Context) func(c *loggerPointConfig)

// default option
func withTraceId() Option {
	return func(ctx context.Context) func(c *loggerPointConfig) {
		return func(c *loggerPointConfig) {
			c.field = append(c.field, zap.String("trace_id", getTraceId(ctx)))
		}
	}
}

func (c *CoreLogger) WithContext(ctx context.Context, opts ...Option) *helper {

	cfg := &loggerPointConfig{}
	opts = append(opts, withTraceId())

	for _, v := range opts {
		v(ctx)(cfg)
	}

	return &helper{logger: &loggerPoint{
		logger: c.logger,
		pool:   c.pool,
		field:  cfg.field,
		ctx:    ctx,
	}}
}

type loggerPointConfig struct {
	field []zap.Field
}

type loggerPoint struct {
	ctx    context.Context
	pool   *sync.Pool
	logger *zap.Logger
	field  []zap.Field
}

func (c *loggerPoint) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}

	buf := c.pool.Get().(*bytes.Buffer)

	for i := 0; i < len(keyvals); i += 2 {
		fmt.Fprintf(buf, "  %s=%v", keyvals[i], keyvals[i+1])
	}

	switch level {
	case log.LevelDebug:
		c.logger.Debug(buf.String(), c.field...)
	case log.LevelInfo:
		c.logger.Info(buf.String(), c.field...)
	case log.LevelWarn:
		c.logger.Warn(buf.String(), c.field...)
	case log.LevelError:
		c.logger.Error(buf.String(), c.field...)
	//default:
	//	c.logger.Debug(buf.String(), c.field...)

	}

	buf.Reset()
	c.pool.Put(buf)
	return nil
}

// get trace id
func getTraceId(ctx context.Context) string {
	var traceID string
	if tid := trace.SpanContextFromContext(ctx).TraceID(); tid.IsValid() {
		traceID = tid.String()
	}
	return traceID
}
