package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logger *Logger) *zap.Logger {
	encoder := newZapEncoder(logger.Format)

	writers := make(map[zapcore.Level]*Rotater)
	if logger.MultiFiles {
		for _, level := range newZapLevels(logger.Level) {
			writers[level] = &Rotater{
				Filename:   logger.Path + "." + level.CapitalString() + ".log",
				MaxSize:    logger.MaxSize,
				MaxBackups: logger.MaxBackups,
				MaxAge:     logger.MaxAge,
				Compress:   logger.Compress,
			}
		}
	} else {
		writers[convertToZapLevel(logger.Level)] = &Rotater{
			Filename:   logger.Path + ".log",
			MaxSize:    logger.MaxSize,
			MaxBackups: logger.MaxBackups,
			MaxAge:     logger.MaxAge,
		}
	}

	cores := make([]zapcore.Core, 0)
	for level, writer := range writers {
		core := zapcore.NewCore(encoder, zapcore.AddSync(writer), newZapLevelEnablerFunc(level, logger.MultiFiles))
		cores = append(cores, core)
	}

	tee := zapcore.NewTee(cores...)

	zlogger := zap.New(tee, zap.AddCaller())
	return zlogger
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
