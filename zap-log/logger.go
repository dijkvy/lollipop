package zap_log

import (
	"fmt"
	"github.com/laxiaohong/lollipop/zap-log/init/config"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// config log level sentinel
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func NewZapLogger(cfg *config.ZapLoggerConfig) (logger *zap.Logger) {
	const (
		_defaultBackUp = 200  // 保留日志的最大值
		_defaultSize   = 1024 // 默认日志最大分割容量
		_defaultAge    = 7    // 日志保留的最大天数
	)

	if cfg == nil {
		panic(fmt.Sprintf("NewZapLogger could't be nil"))
	}

	if cfg.Path == "" {
		panic(fmt.Sprintf("path cound not be null"))
	}

	var err error
	if err = os.MkdirAll(cfg.GetPath(), os.ModePerm); err != nil {
		if os.IsPermission(err) {
			panic(err)
		}
	}

	var logLevel = zap.DebugLevel

	switch cfg.GetLevel() {
	case LevelDebug:
		logLevel = zap.DebugLevel

	case LevelInfo:
		logLevel = zap.InfoLevel
	case LevelWarn:
		logLevel = zap.WarnLevel
	case LevelError:
		logLevel = zap.ErrorLevel
	default:
		panic(fmt.Sprintf("level must be debug, info, warn or error"))

	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		CallerKey:     "file",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		StacktraceKey: "stack",

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.9999999"))
		}, // time format
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 10e6)
		}, //
	}

	// 控制台输出配置
	consoleEncoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		CallerKey:     "file",
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		StacktraceKey: "stack",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.9999999"))
		}, // time format
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 10e6)
		}, //
	}
	// 将所有的日志文件输出到同一个文件
	bizLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	// 保留文件的最大数量
	var maxBackupSize = _defaultBackUp
	if cfg.MaxBackup != nil {
		maxBackupSize = cast.ToInt(cfg.GetMaxBackup())
	}

	// 保留日志的最大天数
	var maxAge = _defaultAge
	if cfg.MaxAge != nil {
		maxAge = cast.ToInt(cfg.GetMaxAge())
	}

	// 日志的最大值
	var maxSize = _defaultSize
	if cfg.MaxSize != nil {
		maxSize = cast.ToInt(cfg.GetMaxSize())
	}

	// writer
	bizWriter := getWriter(cfg.Path+string(filepath.Separator)+"biz.log", maxBackupSize, maxAge, maxSize, cfg.GetCompress())

	// 输出多个
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(bizWriter), bizLevel), // biz 日志
	)

	// debug 日志级别是否输出到控制台
	if cfg.DebugModeOutputConsole != nil && (*cfg.DebugModeOutputConsole && (strings.ToLower(cfg.Level) == "debug")) {

		core = zapcore.NewTee(core,
			zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), zapcore.AddSync(os.Stdout), logLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		)
	}

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

func getWriter(filename string, maxBackup, maxAge, maxSize int, compress bool) io.Writer {

	fmt.Printf("getWriter %s, maxBackup:%d, maxAge:%d, maxSize:%dmb, compress:%v\n", filename, maxBackup, maxAge, maxSize, compress)

	return &lumberjack.Logger{
		Filename:   filename,  // 文件名
		MaxSize:    maxSize,   // the file max size, unit is mb, if overflow the file will rotate
		MaxBackups: maxBackup, // 最大文件保留数, 超过就删除最老的日志文件
		MaxAge:     maxAge,    // 保留文件的最大天数
		Compress:   compress,  // 不启用压缩的功能
		LocalTime:  true,      // 本地时间分割
	}
}
