package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/example/addpay-go/types"
)

// Level represents log levels
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// DefaultLogger is a simple implementation of the Logger interface
type DefaultLogger struct {
	level  Level
	logger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(level Level) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...types.Field) {
	if l.level <= DEBUG {
		l.log("DEBUG", msg, fields...)
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...types.Field) {
	if l.level <= INFO {
		l.log("INFO", msg, fields...)
	}
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...types.Field) {
	if l.level <= WARN {
		l.log("WARN", msg, fields...)
	}
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...types.Field) {
	if l.level <= ERROR {
		l.log("ERROR", msg, fields...)
	}
}

// log formats and logs a message with fields
func (l *DefaultLogger) log(level, msg string, fields ...types.Field) {
	timestamp := time.Now().Format(time.RFC3339)
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level, msg)

	if len(fields) > 0 {
		logMsg += " |"
		for _, field := range fields {
			logMsg += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	l.logger.Println(logMsg)
}

// NoOpLogger is a logger that does nothing
type NoOpLogger struct{}

// NewNoOpLogger creates a new no-op logger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug does nothing
func (l *NoOpLogger) Debug(msg string, fields ...types.Field) {}

// Info does nothing
func (l *NoOpLogger) Info(msg string, fields ...types.Field) {}

// Warn does nothing
func (l *NoOpLogger) Warn(msg string, fields ...types.Field) {}

// Error does nothing
func (l *NoOpLogger) Error(msg string, fields ...types.Field) {}
