package zerolog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/supchat-lmrt/back-go/internal/logger"
)

// ZerologLogger is a wrapper around zerolog.Logger that implements the Logger interface.
type ZerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a new ZerologLogger instance.
func NewZerologLogger(options ...logger.CreateLoggerOption) func() logger.Logger {
	opts := &logger.CreateLoggerOptions{
		MinLevel: logger.LogLevelInfo, // Default minimum log level
	}

	for _, opt := range options {
		opt(opts)
	}

	// colorMagenta = iota + 30 + 1 * 5 (https://en.wikipedia.org/wiki/ANSI_escape_code#Colors)
	zerolog.LevelColors[zerolog.DebugLevel] = 35

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01/2006 - 15:04:05"}

	skipFrameCount := 3 // Skip the zerolog package frames
	logg := &ZerologLogger{
		logger: zerolog.New(output).
			Level(zerolog.Level(opts.MinLevel)).
			With().
			CallerWithSkipFrameCount(skipFrameCount).
			Timestamp().
			Logger(),
	}

	return func() logger.Logger {
		return logg
	}
}

// With returns a new LogEntry to create a nested logger.
func (z *ZerologLogger) With() logger.LogEntry {
	ctx := z.logger.With()
	return &ZerologLogEntry{ctx: &ctx}
}

func (z *ZerologLogger) Trace() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Trace()}
}

// Debug starts a new log entry at Debug level.
func (z *ZerologLogger) Debug() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Debug()}
}

// Info starts a new log entry at Info level.
func (z *ZerologLogger) Info() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Info()}
}

// Warn starts a new log entry at Warn level.
func (z *ZerologLogger) Warn() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Warn()}
}

// Error starts a new log entry at Error level.
func (z *ZerologLogger) Error() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Error()}
}

// Fatal starts a new log entry at Fatal level.
func (z *ZerologLogger) Fatal() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Fatal()}
}

// Panic starts a new log entry at Panic level.
func (z *ZerologLogger) Panic() logger.LogEntry {
	return &ZerologLogEntry{event: z.logger.Panic()}
}

// ZerologLogEntry is a wrapper around zerolog.Event or zerolog.Context that implements the LogEntry interface.
type ZerologLogEntry struct {
	ctx   *zerolog.Context // Used for With()
	event *zerolog.Event   // Used for Debug(), Info(), etc.
}

// Str adds a string field to the log entry.
func (z *ZerologLogEntry) Str(key, value string) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Str(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Str(key, value)
	}
	return z
}

// Bool adds a boolean field to the log entry.
func (z *ZerologLogEntry) Bool(key string, value bool) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Bool(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Bool(key, value)
	}
	return z
}

// Int adds an integer field to the log entry.
func (z *ZerologLogEntry) Int(key string, value int) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Int(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Int(key, value)
	}
	return z
}

// Float64 adds a float64 field to the log entry.
func (z *ZerologLogEntry) Float64(key string, value float64) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Float64(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Float64(key, value)
	}
	return z
}

// Dur adds a time.Duration field to the log entry.
func (z *ZerologLogEntry) Dur(key string, value time.Duration) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Dur(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Dur(key, value)
	}
	return z
}

// Any adds any field to the log entry.
func (z *ZerologLogEntry) Any(key string, value any) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Interface(key, value)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Interface(key, value)
	}
	return z
}

func (z *ZerologLogEntry) Err(err error) logger.LogEntry {
	if z.ctx != nil {
		ctx := z.ctx.Err(err)
		z.ctx = &ctx
	} else if z.event != nil {
		z.event = z.event.Err(err)
	}
	return z
}

// Msg sends the log entry with the given message.
func (z *ZerologLogEntry) Msg(msg string) {
	if z.event != nil {
		z.event.Msg(msg)
	}
}

func (z *ZerologLogEntry) Msgf(format string, v ...any) {
	if z.event != nil {
		z.event.Msgf(format, v...)
	}
}

func (z *ZerologLogEntry) Send() {
	if z.event != nil {
		z.event.Send()
	}
}

// Logger returns a new logger with the fields added to the log entry.
func (z *ZerologLogEntry) Logger() logger.Logger {
	if z.ctx != nil {
		return &ZerologLogger{logger: z.ctx.Logger()}
	}
	return &ZerologLogger{logger: zerolog.New(os.Stderr).With().Timestamp().Logger()}
}
