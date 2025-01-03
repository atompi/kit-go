package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(opts *Options) *zap.Logger {
	encoder := newZapEncoder(opts.Format)

	writers := make(map[zapcore.Level]*Logger)
	if opts.MultiFiles {
		for _, level := range newZapLevels(opts.Level) {
			writers[level] = &Logger{
				Filename:   opts.Path + "." + level.CapitalString() + ".log",
				MaxSize:    opts.MaxSize,
				MaxBackups: opts.MaxBackups,
				MaxAge:     opts.MaxAge,
				Compress:   opts.Compress,
			}
		}
	} else {
		writers[convertToZapLevel(opts.Level)] = &Logger{
			Filename:   opts.Path + ".log",
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
		}
	}

	cores := make([]zapcore.Core, 0)
	for level, writer := range writers {
		core := zapcore.NewCore(encoder, zapcore.AddSync(writer), newZapLevelEnablerFunc(level, opts.MultiFiles))
		cores = append(cores, core)
	}

	tee := zapcore.NewTee(cores...)

	logger := zap.New(tee, zap.AddCaller())
	return logger
}

func newZapEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	switch format {
	case "json", "JSON":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func newZapLevelEnablerFunc(l zapcore.Level, m bool) zap.LevelEnablerFunc {
	switch l {
	case zap.DebugLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			if !m {
				return level >= zapcore.DebugLevel
			}
			return level == zapcore.DebugLevel
		})
	case zap.InfoLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			if !m {
				return level >= zapcore.InfoLevel
			}
			return level == zapcore.InfoLevel
		})
	case zap.WarnLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			if !m {
				return level >= zapcore.WarnLevel
			}
			return level == zapcore.WarnLevel
		})
	case zap.ErrorLevel:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			if !m {
				return level >= zapcore.ErrorLevel
			}
			return level == zapcore.ErrorLevel
		})
	default:
		return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			if !m {
				return level >= zapcore.DebugLevel
			}
			return level == zapcore.DebugLevel
		})
	}
}

func convertToZapLevel(level string) zapcore.Level {
	switch level {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO":
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

func newZapLevels(level string) []zapcore.Level {
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}

	switch level {
	case "debug", "DEBUG":
		levels = levels[:]
	case "info", "INFO":
		levels = levels[1:]
	case "warn", "WARN":
		levels = levels[2:]
	case "error", "ERROR":
		levels = levels[3:]
	default:
	}

	return levels
}
