package log

type Logger struct {
	Level      string
	Format     string
	Path       string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	MultiFiles bool
}

type Options func(*Logger)

func NewLoggerOptions(opts ...Options) *Logger {
	logger := &Logger{
		Level:      "info",
		Format:     "console",
		Path:       "logger",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
		MultiFiles: false,
	}
	for _, f := range opts {
		f(logger)
	}
	return logger
}

func WithLevel(level string) Options {
	return func(logger *Logger) {
		logger.Level = level
	}
}

func WithFormat(format string) Options {
	return func(logger *Logger) {
		logger.Format = format
	}
}

func WithPath(path string) Options {
	return func(logger *Logger) {
		logger.Path = path
	}
}

func WithMaxSize(maxSize int) Options {
	return func(logger *Logger) {
		logger.MaxSize = maxSize
	}
}

func WithMaxAge(maxAge int) Options {
	return func(logger *Logger) {
		logger.MaxAge = maxAge
	}
}

func WithMaxBackups(maxBackups int) Options {
	return func(logger *Logger) {
		logger.MaxBackups = maxBackups
	}
}

func WithCompress(compress bool) Options {
	return func(logger *Logger) {
		logger.Compress = compress
	}
}

func WithMultiFiles(multiFiles bool) Options {
	return func(logger *Logger) {
		logger.MultiFiles = multiFiles
	}
}
