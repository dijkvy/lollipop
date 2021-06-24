package log

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"sync"
)

func NewKratosLog(logger *zap.Logger) log.Logger {
	return &kratosLog{
		logger: logger,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

type kratosLog struct {
	logger *zap.Logger
	pool   *sync.Pool
}

func (c *kratosLog) Log(level log.Level, keyvals ...interface{}) error {
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
		c.logger.Debug(buf.String())
	case log.LevelInfo:
		c.logger.Info(buf.String())
	case log.LevelWarn:
		c.logger.Warn(buf.String())
	case log.LevelError:
		c.logger.Error(buf.String())
	}

	buf.Reset()
	c.pool.Put(buf)
	return nil
}
