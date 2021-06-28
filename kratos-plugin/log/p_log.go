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

func NewCoreLogger(logger *zap.Logger, opts ...Option) *CoreLogger {
	return &CoreLogger{
		logger: logger,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		option: []Option{withTraceId()},
	}
}

type CoreLogger struct {
	logger *zap.Logger
	pool   *sync.Pool
	option []Option
}

// Option type
type Option = func(ctx context.Context) zap.Field

// default option
func withTraceId() Option {
	return func(ctx context.Context) (field zap.Field) {
		return zap.String("trace_id", getTraceId(ctx))
	}
}

func (c *CoreLogger) WithContext(ctx context.Context) *helper {

	var field []zap.Field
	for _, v := range c.option {
		field = append(field, v(ctx))
	}

	return &helper{logger: &loggerPoint{
		logger: c.logger,
		pool:   c.pool,
		field:  field,
		ctx:    ctx,
	}}
}

type LoggerPointConfig struct {
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
