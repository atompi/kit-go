package log

import (
	"context"
	"log/slog"
	"path/filepath"
)

type handler struct {
	level       slog.Level
	writers     map[slog.Level]*Logger
	options     *Options
	slogOptions *slog.HandlerOptions
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	writer := &Logger{}
	if !h.options.MultiFiles {
		writer = h.writers[convertToSlogLevel(h.options.Level)]
	} else {
		writer = h.writers[r.Level]
	}

	switch h.options.Format {
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

func newSlogHandler(opts *Options, slogOpts *slog.HandlerOptions) *handler {
	h := &handler{
		level:       convertToSlogLevel(opts.Level),
		writers:     make(map[slog.Level]*Logger),
		options:     opts,
		slogOptions: slogOpts,
	}

	if opts.MultiFiles {
		for _, level := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
			h.writers[level] = &Logger{
				Filename:   opts.Path + "." + level.String() + ".log",
				MaxSize:    opts.MaxSize,
				MaxBackups: opts.MaxBackups,
				MaxAge:     opts.MaxAge,
				Compress:   opts.Compress,
			}
		}
	} else {
		h.writers[convertToSlogLevel(opts.Level)] = &Logger{
			Filename:   opts.Path + ".log",
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}
	}

	return h
}

func NewSlogLogger(opts *Options) *slog.Logger {
	slogOpts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       convertToSlogLevel(opts.Level),
		ReplaceAttr: newSlogReplaceAttr(),
	}

	logger := slog.New(newSlogHandler(opts, slogOpts))
	return logger
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
