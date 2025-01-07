package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Config holds logger configuration
type Config struct {
	Level      string
	Debug      bool
	TimeFormat string
	Pretty     bool
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg Config) {
	// Set default time format if not specified
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}

	// Configure logger output
	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.TimeFormat,
		}
	}

	// Set global logger
	zerolog.TimeFieldFormat = cfg.TimeFormat
	level := getLogLevel(cfg.Level)
	if cfg.Debug {
		level = zerolog.DebugLevel
	}

	log = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}

// GetLogger returns the configured logger instance
func GetLogger() zerolog.Logger {
	return log
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return log.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// WithError adds an error to the log event
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// WithField adds a field to the log event
func WithField(key string, value interface{}) zerolog.Logger {
	return log.With().Interface(key, value).Logger()
}

// getLogLevel converts a string level to zerolog.Level
func getLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
