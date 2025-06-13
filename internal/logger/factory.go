package logger

type CreateLoggerOptions struct {
	MinLevel LogLevel
}

type CreateLoggerOption func(*CreateLoggerOptions)

func WithMinLevel(level LogLevel) CreateLoggerOption {
	return func(opts *CreateLoggerOptions) {
		opts.MinLevel = level
	}
}
