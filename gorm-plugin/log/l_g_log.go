package log

// gorm log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"sync"
	"time"
)

var (
	tracer = otel.Tracer("github.com/laxiaohong/x")

	pool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

func NewLogger(log *zap.Logger, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s[info] "
		warnStr      = "%s[warn] "
		errStr       = "%s[error] "
		traceStr     = "%s[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &loggerPoint{
		logger:       log,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type loggerPoint struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
	logger                              *zap.Logger
}

func (l *loggerPoint) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l loggerPoint) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l._info(ctx, l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l loggerPoint) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l._warn(ctx, l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l loggerPoint) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l._error(ctx, l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l loggerPoint) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l._error(ctx, l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l._error(ctx, l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l._warn(ctx, l.traceErrStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l._warn(ctx, l.traceErrStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l._info(ctx, l.traceErrStr, utils.FileWithLineNum(), utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l._info(ctx, l.traceErrStr, utils.FileWithLineNum(), utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (l *loggerPoint) _info(ctx context.Context, msg string, data ...interface{}) {
	buff := pool.Get().(*bytes.Buffer)
	fmt.Fprintf(buff, msg, data...)
	l.logger.Info(buff.String(), zap.String("trace_id", getTraceId(ctx)))
	buff.Reset()
	pool.Put(buff)
}

func (l *loggerPoint) _warn(ctx context.Context, msg string, data ...interface{}) {
	buff := pool.Get().(*bytes.Buffer)
	fmt.Fprintf(buff, msg, data...)
	l.logger.Warn(buff.String(), zap.String("trace_id", getTraceId(ctx)))
	buff.Reset()
	pool.Put(buff)
}

func (l *loggerPoint) _error(ctx context.Context, msg string, data ...interface{}) {
	buff := pool.Get().(*bytes.Buffer)
	fmt.Fprintf(buff, msg, data...)
	l.logger.Error(buff.String(), zap.String("trace_id", getTraceId(ctx)))
	buff.Reset()
	pool.Put(buff)

}

// get trace id
func getTraceId(ctx context.Context) string {
	var traceID string
	if tid := trace.SpanContextFromContext(ctx).TraceID(); tid.IsValid() {
		traceID = tid.String()
	}
	return traceID
}
