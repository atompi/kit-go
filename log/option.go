package log

type Options struct {
	Level      string
	Format     string
	Path       string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	MultiFiles bool
}

func NewOptions() *Options {
	return &Options{
		Level:      "info",
		Format:     "console",
		Path:       "logger",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
		MultiFiles: false,
	}
}
