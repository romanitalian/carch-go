package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is a wrapper around zerolog.Logger
type Logger struct {
	logger zerolog.Logger
}

// Option is a function that configures a Logger
type Option func(*Logger)

// New creates a new logger with the given options
func New(options ...Option) *Logger {
	// Default configuration
	l := &Logger{
		logger: zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(zerolog.InfoLevel),
	}

	// Apply options
	for _, option := range options {
		option(l)
	}

	return l
}

// WithLevel sets the logger level
func WithLevel(level zerolog.Level) Option {
	return func(l *Logger) {
		l.logger = l.logger.Level(level)
	}
}

// WithOutput sets the logger output
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.logger = zerolog.New(w).
			With().
			Timestamp().
			Caller().
			Logger()
	}
}

// WithPretty enables pretty logging
func WithPretty() Option {
	return func(l *Logger) {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		l.logger = zerolog.New(output).
			With().
			Timestamp().
			Caller().
			Logger()
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// GetZerologLogger returns the underlying zerolog.Logger
func (l *Logger) GetZerologLogger() zerolog.Logger {
	return l.logger
}
