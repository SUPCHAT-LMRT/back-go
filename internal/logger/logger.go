package logger

import (
	"time"
)

// Logger is the interface that all loggers should implement.
type Logger interface {
	// With returns a new logger with the given fields added to the existing ones.
	With() LogEntry

	// Trace starts a new log entry at Trace level.
	Trace() LogEntry

	// Debug starts a new log entry at Debug level.
	Debug() LogEntry

	// Info starts a new log entry at Info level.
	Info() LogEntry

	// Warn starts a new log entry at Warn level.
	Warn() LogEntry

	// Error starts a new log entry at Error level.
	Error() LogEntry

	// Fatal starts a new log entry at Fatal level.
	Fatal() LogEntry

	// Panic starts a new log entry at Panic level.
	Panic() LogEntry
}

// LogEntry represents a log entry with chained methods.
type LogEntry interface {
	// Str adds a string field to the log entry.
	Str(key, value string) LogEntry

	// Bool adds a boolean field to the log entry.
	Bool(key string, value bool) LogEntry

	// Int adds an integer field to the log entry.
	Int(key string, value int) LogEntry

	// Float64 adds a float64 field to the log entry.
	Float64(key string, value float64) LogEntry

	// Dur adds a time.Duration field to the log entry.
	Dur(key string, value time.Duration) LogEntry

	// Any adds any field to the log entry.
	Any(key string, value interface{}) LogEntry

	// Err adds an error field to the log entry.
	Err(err error) LogEntry

	// Msg sends the log entry with the given message.
	Msg(msg string)
	// Msgf sends the log entry with the given message.
	Msgf(format string, v ...interface{})

	// Send sends the log entry.
	Send()

	// Logger returns a new logger with the fields added to the log entry.
	Logger() Logger
}
