package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option struct {
	LogLevel   string
	Format     string
	LogPath    string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

func NewOptions() *Option {
	return &Option{
		LogLevel:   "info",
		Format:     "console",
		LogPath:    "logger",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
	}
}

func InitLogger(opt *Option) *zap.Logger {
	encoder := newEncoder(opt.Format)

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, newLogWriter("debug", opt), newLevelEnablerFunc(zapcore.DebugLevel)),
		zapcore.NewCore(encoder, newLogWriter("info", opt), newLevelEnablerFunc(zapcore.InfoLevel)),
		zapcore.NewCore(encoder, newLogWriter("warn", opt), newLevelEnablerFunc(zapcore.WarnLevel)),
		zapcore.NewCore(encoder, newLogWriter("error", opt), newLevelEnablerFunc(zapcore.ErrorLevel)),
	}

	switch opt.LogLevel {
	case "debug", "DEBUG":
		cores = cores[:]
	case "info", "INFO":
		cores = cores[1:]
	case "warn", "WARN":
		cores = cores[2:]
	case "error", "ERROR":
		cores = cores[3:]
	default:
	}

	tee := zapcore.NewTee(cores...)

	logger := zap.New(tee, zap.AddCaller())
	return logger
}

func newEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	switch format {
	case "json", "JSON":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console", "text", "TEXT":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func newLogWriter(level string, opt *Option) zapcore.WriteSyncer {
	filename := opt.LogPath + "." + level + ".log"

	logger := &Logger{
		Filename:   filename,
		MaxSize:    opt.MaxSize,
		MaxAge:     opt.MaxAge,
		MaxBackups: opt.MaxBackups,
		Compress:   opt.Compress,
	}

	return zapcore.AddSync(logger)
}

func newLevelEnablerFunc(l zapcore.Level) zap.LevelEnablerFunc {
	switch l {
	case zap.DebugLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level == zapcore.DebugLevel
		})
	case zap.InfoLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level == zapcore.InfoLevel
		})
	case zap.WarnLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level == zapcore.WarnLevel
		})
	case zap.ErrorLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level == zapcore.ErrorLevel
		})
	default:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return false
		})
	}
}
