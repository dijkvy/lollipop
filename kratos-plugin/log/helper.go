package log

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type helper struct {
	logger log.Logger
}

// Debug logs a message at debug level.
func (h *helper) Debug(a ...interface{}) {

	h.logger.Log(log.LevelDebug, "msg", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *helper) Debugf(format string, a ...interface{}) {
	h.logger.Log(log.LevelDebug, "msg", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *helper) Debugw(keyvals ...interface{}) {
	h.logger.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *helper) Info(a ...interface{}) {
	h.logger.Log(log.LevelInfo, "msg", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *helper) Infof(format string, a ...interface{}) {
	h.logger.Log(log.LevelInfo, "msg", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *helper) Infow(keyvals ...interface{}) {
	h.logger.Log(log.LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *helper) Warn(a ...interface{}) {
	h.logger.Log(log.LevelWarn, "msg", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *helper) Warnf(format string, a ...interface{}) {
	h.logger.Log(log.LevelWarn, "msg", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *helper) Warnw(keyvals ...interface{}) {
	h.logger.Log(log.LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *helper) Error(a ...interface{}) {
	h.logger.Log(log.LevelError, "msg", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *helper) Errorf(format string, a ...interface{}) {
	h.logger.Log(log.LevelError, "msg", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *helper) Errorw(keyvals ...interface{}) {
	h.logger.Log(log.LevelError, keyvals...)
}
