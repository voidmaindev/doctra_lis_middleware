// Package log provides a wrapper around the zerolog package to log messages at different levels.
package log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/voidmaindev/doctra_lis_middleware/config"
)

// Logger is a wrapper around zerolog.Logger
type Logger struct {
	Logger   *zerolog.Logger
	Disabled bool
}

// NewLogger creates a new Logger instance
func NewLogger() (*Logger, error) {
	settings, err := config.ReadLogConfig()
	if err != nil {
		return nil, err
	}

	if settings.Disable {
		return &Logger{Disabled: true}, nil
	}

	cw := zerolog.ConsoleWriter{
		Out:        getOutput(settings.Output),
		TimeFormat: getTimestampFormat(settings.TimeFormat),
	}

	loggerContext := zerolog.New(cw).With().Timestamp()

	if settings.AddCaller {
		loggerContext = loggerContext.Caller()
	}

	if settings.AddPid {
		loggerContext = loggerContext.Int("pid", os.Getpid())
	}

	logger := loggerContext.Logger()

	return &Logger{Logger: &logger}, nil
}

// Trace logs a message at trace level
func (l *Logger) Trace(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Trace().Msg(msg)
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Debug().Msg(msg)
}

// Info logs a message at info level
func (l *Logger) Info(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Info().Msg(msg)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Warn().Msg(msg)
}

// Error logs a message at error level
func (l *Logger) Error(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Error().Msg(msg)
}

// Err logs a message at error level with an error
func (l *Logger) Err(err error, msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Err(err).Msg(msg)
}

// Panic logs a message at panic level
func (l *Logger) Panic(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Panic().Msg(msg)
}

// Fatal logs a message at fatal level
func (l *Logger) Fatal(msg string) {
	if l.Disabled {
		return
	}

	(*l.Logger).Fatal().Msg(msg)
}

// getOutput returns the output writer based on the output string
func getOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return os.Stderr
	}
}

// getTimestampFormat returns the timestamp format based on the timeFormat string
func getTimestampFormat(timeFormat string) string {
	switch strings.ToLower(timeFormat) {
	case "rfc3339":
		return time.RFC3339
	case "unix":
		return zerolog.TimeFormatUnix
	case "unix_ms":
		return zerolog.TimeFormatUnixMs
	case "unix_us":
		return zerolog.TimeFormatUnixMicro
	case "unix_ns":
		return zerolog.TimeFormatUnixNano
	default:
		return zerolog.TimeFormatUnix
	}
}
