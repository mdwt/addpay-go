package logger

import (
	"log/slog"
	"os"

	"github.com/mdwt/addpay-go/types"
)

// SlogLogger wraps slog.Logger to implement the simple Logger interface
type SlogLogger struct {
	logger *slog.Logger
}

// NewDefaultLogger creates a default slog-based logger
func NewDefaultLogger() types.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return &SlogLogger{logger: logger}
}

// Debug logs a debug message with key-value pairs
func (s *SlogLogger) Debug(msg string, keysAndValues ...interface{}) {
	s.logger.Debug(msg, keysAndValues...)
}

// Info logs an info message with key-value pairs
func (s *SlogLogger) Info(msg string, keysAndValues ...interface{}) {
	s.logger.Info(msg, keysAndValues...)
}

// Warn logs a warning message with key-value pairs
func (s *SlogLogger) Warn(msg string, keysAndValues ...interface{}) {
	s.logger.Warn(msg, keysAndValues...)
}

// Error logs an error message with key-value pairs
func (s *SlogLogger) Error(msg string, keysAndValues ...interface{}) {
	s.logger.Error(msg, keysAndValues...)
}

// NoOpLogger discards all log messages
type NoOpLogger struct{}

// NewNoOpLogger creates a logger that discards all messages
func NewNoOpLogger() types.Logger {
	return &NoOpLogger{}
}

// Debug does nothing
func (n *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Info does nothing
func (n *NoOpLogger) Info(msg string, keysAndValues ...interface{}) {}

// Warn does nothing
func (n *NoOpLogger) Warn(msg string, keysAndValues ...interface{}) {}

// Error does nothing
func (n *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {}
