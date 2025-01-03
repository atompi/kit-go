package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(opts *Options) *zap.Logger {
	encoder := newEncoder(opts.Format)

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, newLogWriter("DEBUG", opts), newLevelEnablerFunc(zapcore.DebugLevel)),
		zapcore.NewCore(encoder, newLogWriter("INFO", opts), newLevelEnablerFunc(zapcore.InfoLevel)),
		zapcore.NewCore(encoder, newLogWriter("WARN", opts), newLevelEnablerFunc(zapcore.WarnLevel)),
		zapcore.NewCore(encoder, newLogWriter("ERROR", opts), newLevelEnablerFunc(zapcore.ErrorLevel)),
	}

	switch opts.Level {
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

func newLogWriter(level string, opts *Options) zapcore.WriteSyncer {
	filename := opts.Path + "." + level + ".log"

	logger := &Logger{
		Filename:   filename,
		MaxSize:    opts.MaxSize,
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxBackups,
		Compress:   opts.Compress,
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
