package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logLevel string, logPath string, maxSize int, maxAge int, compress bool) *zap.Logger {
	debugWriteSyncer := newLogWriter(logPath+".debug.log", maxSize, maxAge, compress)
	infoWriteSyncer := newLogWriter(logPath+".info.log", maxSize, maxAge, compress)
	warnWriteSyncer := newLogWriter(logPath+".warn.log", maxSize, maxAge, compress)
	errorWriteSyncer := newLogWriter(logPath+".error.log", maxSize, maxAge, compress)

	encoder := newEncoder()

	debugCore := zapcore.NewCore(encoder, debugWriteSyncer, newLevelEnablerFunc(zapcore.DebugLevel))
	infoCore := zapcore.NewCore(encoder, infoWriteSyncer, newLevelEnablerFunc(zapcore.InfoLevel))
	warnCore := zapcore.NewCore(encoder, warnWriteSyncer, newLevelEnablerFunc(zapcore.WarnLevel))
	errorCore := zapcore.NewCore(encoder, errorWriteSyncer, newLevelEnablerFunc(zapcore.ErrorLevel))

	tee := zapcore.NewTee(debugCore, infoCore, warnCore, errorCore)

	logger := zap.New(tee, zap.AddCaller())
	return logger
}

func newEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newLogWriter(logPath string, maxSize int, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: logPath,
		MaxSize:  maxSize,
		MaxAge:   maxAge,
		Compress: compress,
	}
	return zapcore.AddSync(lumberJackLogger)
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
