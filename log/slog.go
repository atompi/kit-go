package log

import (
	"context"
	"log/slog"
	"path/filepath"
)

type handler struct {
	level       slog.Level
	writers     map[slog.Level]*Rotater
	logger      *Logger
	slogOptions *slog.HandlerOptions
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	writer := &Rotater{}
	if !h.logger.MultiFiles {
		writer = h.writers[convertToSlogLevel(h.logger.Level)]
	} else {
		writer = h.writers[r.Level]
	}

	switch h.logger.Format {
	case "json", "JSON":
		return slog.NewJSONHandler(writer, h.slogOptions).Handle(ctx, r)
	default:
		return slog.NewTextHandler(writer, h.slogOptions).Handle(ctx, r)
	}
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h
}

func newSlogHandler(logger *Logger, slogOptions *slog.HandlerOptions) *handler {
	h := &handler{
		level:       convertToSlogLevel(logger.Level),
		writers:     make(map[slog.Level]*Rotater),
		logger:      logger,
		slogOptions: slogOptions,
	}

	if logger.MultiFiles {
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			h.writers[level] = &Rotater{
				Filename:   logger.Path + "." + level.String() + ".log",
				MaxSize:    logger.MaxSize,
				MaxBackups: logger.MaxBackups,
				MaxAge:     logger.MaxAge,
				Compress:   logger.Compress,
			}
		}
	} else {
		h.writers[convertToSlogLevel(logger.Level)] = &Rotater{
			Filename:   logger.Path + ".log",
			MaxSize:    logger.MaxSize,
			MaxBackups: logger.MaxBackups,
			MaxAge:     logger.MaxAge,
			Compress:   logger.Compress,
		}
	}

	return h
}

func NewSlogLogger(logger *Logger) *slog.Logger {
	slogOptions := &slog.HandlerOptions{
		AddSource:   true,
		Level:       convertToSlogLevel(logger.Level),
		ReplaceAttr: newSlogReplaceAttr(),
	}

	slogger := slog.New(newSlogHandler(logger, slogOptions))
	return slogger
}

func convertToSlogLevel(level string) slog.Level {
	switch level {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func newSlogReplaceAttr() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}
			if source != nil {
				baseDir, err := filepath.Abs(".")
				if err != nil {
					baseDir = filepath.Dir(source.File)
				}
				relPath, err := filepath.Rel(baseDir, source.File)
				if err != nil {
					relPath = filepath.Base(source.File)
				}
				source.File = relPath
			}
		}
		return a
	}
}
