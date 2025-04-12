package executor

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel defines the level of logging
type LogLevel int

const (
	// LogLevelDebug includes all messages
	LogLevelDebug LogLevel = iota
	// LogLevelInfo includes info, warn, error, and fatal messages
	LogLevelInfo
	// LogLevelWarn includes warn, error, and fatal messages
	LogLevelWarn
	// LogLevelError includes error and fatal messages
	LogLevelError
	// LogLevelFatal includes only fatal messages
	LogLevelFatal
	// LogLevelNone disables all logging
	LogLevelNone
)

// Logger defines the interface for logging
type Logger interface {
	// Debugf logs a debug message
	Debugf(format string, args ...interface{})
	// Infof logs an info message
	Infof(format string, args ...interface{})
	// Warnf logs a warning message
	Warnf(format string, args ...interface{})
	// Errorf logs an error message
	Errorf(format string, args ...interface{})
	// Fatalf logs a fatal message and exits
	Fatalf(format string, args ...interface{})
	// SetLevel sets the log level
	SetLevel(level LogLevel)
	// GetLevel gets the current log level
	GetLevel() LogLevel
}

// defaultLogger implements Logger with standard log package
type defaultLogger struct {
	level  LogLevel
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
	fatal  *log.Logger
}

// newDefaultLogger creates a new default logger
func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		level:  LogLevelInfo, // Default log level
		debug:  log.New(os.Stdout, "[DEBUG] ", log.Ltime),
		info:   log.New(os.Stdout, "[INFO] ", log.Ltime),
		warn:   log.New(os.Stdout, "[WARN] ", log.Ltime),
		error:  log.New(os.Stderr, "[ERROR] ", log.Ltime),
		fatal:  log.New(os.Stderr, "[FATAL] ", log.Ltime),
	}
}

// Debugf logs a debug message
func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.debug.Output(2, fmt.Sprintf(format, args...))
	}
}

// Infof logs an info message
func (l *defaultLogger) Infof(format string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.info.Output(2, fmt.Sprintf(format, args...))
	}
}

// Warnf logs a warning message
func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.warn.Output(2, fmt.Sprintf(format, args...))
	}
}

// Errorf logs an error message
func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.error.Output(2, fmt.Sprintf(format, args...))
	}
}

// Fatalf logs a fatal message and exits
func (l *defaultLogger) Fatalf(format string, args ...interface{}) {
	if l.level <= LogLevelFatal {
		l.fatal.Output(2, fmt.Sprintf(format, args...))
		os.Exit(1)
	}
}

// SetLevel sets the log level
func (l *defaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel gets the current log level
func (l *defaultLogger) GetLevel() LogLevel {
	return l.level
}

// nullLogger implements Logger with no output
type nullLogger struct{}

// newNullLogger creates a new null logger
func newNullLogger() *nullLogger {
	return &nullLogger{}
}

// Debugf implements Logger.Debugf
func (l *nullLogger) Debugf(format string, args ...interface{}) {}

// Infof implements Logger.Infof
func (l *nullLogger) Infof(format string, args ...interface{}) {}

// Warnf implements Logger.Warnf
func (l *nullLogger) Warnf(format string, args ...interface{}) {}

// Errorf implements Logger.Errorf
func (l *nullLogger) Errorf(format string, args ...interface{}) {}

// Fatalf implements Logger.Fatalf
func (l *nullLogger) Fatalf(format string, args ...interface{}) {}

// SetLevel implements Logger.SetLevel
func (l *nullLogger) SetLevel(level LogLevel) {}

// GetLevel implements Logger.GetLevel
func (l *nullLogger) GetLevel() LogLevel {
	return LogLevelNone
}
