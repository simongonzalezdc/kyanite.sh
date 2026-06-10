// Package logging provides debug-aware logging utilities for the application.
// It wraps charmbracelet/log with TUI-aware suppression and build-tag debug control.
package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	charmlog "github.com/charmbracelet/log"
	"github.com/kyanite/noise/internal/config"
	errutil "github.com/kyanite/noise/internal/errutil"
)

// LogLevel represents the severity level of log messages
type LogLevel int

const (
	// DEBUG is the debug log level
	DEBUG LogLevel = iota
	// INFO is the info log level
	INFO
	// WARN is the warning log level
	WARN
	// ERROR is the error log level
	ERROR
	// FATAL is the fatal log level
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func toCharmLevel(l LogLevel) charmlog.Level {
	switch l {
	case DEBUG:
		return charmlog.DebugLevel
	case INFO:
		return charmlog.InfoLevel
	case WARN:
		return charmlog.WarnLevel
	case ERROR:
		return charmlog.ErrorLevel
	case FATAL:
		return charmlog.FatalLevel
	default:
		return charmlog.InfoLevel
	}
}

// Logger wraps charmbracelet/log with TUI-aware suppression.
type Logger struct {
	inner   *charmlog.Logger
	level   LogLevel
	tuiMode bool
	logFile *os.File
}

// Config holds logger configuration
type Config struct {
	Level      LogLevel
	Output     io.Writer
	ShowCaller bool
	LogFile    string
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      INFO,
		Output:     os.Stderr, // CRITICAL: Use stderr to avoid corrupting Bubble Tea's stdout rendering
		ShowCaller: false,
		LogFile:    "",
	}
}

// New creates a new logger with the given configuration
func New(cfg *Config) (*Logger, error) {
	output := cfg.Output

	// If log file is specified, create or append to it
	var logFile *os.File
	if cfg.LogFile != "" {
		logDir := filepath.Dir(cfg.LogFile)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return nil, errutil.Wrap(err, "create log directory")
		}

		var err error
		logFile, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, errutil.Wrap(err, "open log file")
		}

		// Use multi-writer to write to both file and stderr
		output = io.MultiWriter(cfg.Output, logFile)
	}

	inner := charmlog.NewWithOptions(output, charmlog.Options{
		ReportTimestamp: true,
		ReportCaller:    cfg.ShowCaller,
		Level:           toCharmLevel(cfg.Level),
	})

	return &Logger{
		inner:   inner,
		level:   cfg.Level,
		logFile: logFile,
	}, nil
}

// Close closes the log file if one was opened
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// NewFromConfig creates a new logger from application configuration
func NewFromConfig(appConfig *config.Config) (*Logger, error) {
	// Parse log level from config
	level := parseLogLevel(appConfig.Dev.LogLevel)

	// Allow forcing debug logging through environment variable or debug build tag
	forcedDebugEnv := os.Getenv("NOISE_DEBUG_LOGGING")
	forcedDebug := forcedDebugEnv != "" &&
		(strings.EqualFold(forcedDebugEnv, "true") || forcedDebugEnv == "1")

	isDebugMode := buildDebugLogging || appConfig.IsDebug() || forcedDebug

	// Ensure debug logging is only active in explicit debug modes
	if isDebugMode {
		level = DEBUG
	} else if level == DEBUG {
		level = INFO
	}

	// Determine output file if in debug mode (runtime debug mode only)
	var logFile string
	if appConfig.IsDebug() {
		logFile = filepath.Join(appConfig.GetDataDir(), "logs", "noise.sh.log")
	}

	cfg := &Config{
		Level:      level,
		Output:     os.Stderr, // CRITICAL: Use stderr to avoid corrupting Bubble Tea's stdout rendering
		ShowCaller: isDebugMode,
		LogFile:    logFile,
	}

	return New(cfg)
}

// parseLogLevel parses a log level string
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
	l.inner.SetLevel(toCharmLevel(level))
}

// SetShowCaller sets whether to show caller information
func (l *Logger) SetShowCaller(show bool) {
	l.inner.SetReportCaller(show)
}

// SetTUIMode enables or disables TUI mode (suppresses output when TUI is active)
func (l *Logger) SetTUIMode(enabled bool) {
	l.tuiMode = enabled
}

// IsTUIMode returns whether TUI mode is active
func (l *Logger) IsTUIMode() bool {
	return l.tuiMode
}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if l.tuiMode || l.level > DEBUG {
		return
	}
	l.inner.Debug(fmt.Sprint(v...))
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.tuiMode || l.level > DEBUG {
		return
	}
	l.inner.Debugf(format, v...)
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	if l.tuiMode || l.level > INFO {
		return
	}
	l.inner.Info(fmt.Sprint(v...))
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.tuiMode || l.level > INFO {
		return
	}
	l.inner.Infof(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(v ...interface{}) {
	if l.tuiMode || l.level > WARN {
		return
	}
	l.inner.Warn(fmt.Sprint(v...))
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.tuiMode || l.level > WARN {
		return
	}
	l.inner.Warnf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.tuiMode || l.level > ERROR {
		return
	}
	l.inner.Error(fmt.Sprint(v...))
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.tuiMode || l.level > ERROR {
		return
	}
	l.inner.Errorf(format, v...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(v ...interface{}) {
	if !l.tuiMode && l.level <= FATAL {
		l.inner.Fatal(fmt.Sprint(v...))
	}
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits the program
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if !l.tuiMode && l.level <= FATAL {
		l.inner.Fatalf(format, v...)
	}
	os.Exit(1)
}

// Global logger instance
var defaultLogger *Logger

// init initializes the default logger
func init() {
	cfg := DefaultConfig()
	logger, err := New(cfg)
	if err != nil {
		// Fallback to a minimal charmbracelet logger if initialization fails
		inner := charmlog.New(os.Stderr)
		inner.SetLevel(charmlog.InfoLevel)
		defaultLogger = &Logger{
			inner: inner,
			level: INFO,
		}
	} else {
		defaultLogger = logger
	}
}

// SetDefaultLogger sets the default logger
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// GetDefaultLogger returns the default logger
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// DebugEnabled reports whether debug-level logging is currently enabled.
func DebugEnabled() bool {
	if defaultLogger == nil {
		return false
	}
	return defaultLogger.level <= DEBUG
}

// EnableTUIMode enables TUI mode which suppresses log output to avoid corrupting the terminal
func EnableTUIMode() {
	if defaultLogger != nil {
		defaultLogger.SetTUIMode(true)
	}
}

// DisableTUIMode disables TUI mode and restores normal log output
func DisableTUIMode() {
	if defaultLogger != nil {
		defaultLogger.SetTUIMode(false)
	}
}

// Convenience functions that use the default logger

// Debug logs a debug message using the default logger
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Debugf logs a formatted debug message using the default logger
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Info logs an info message using the default logger
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infof logs a formatted info message using the default logger
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Warn logs a warning message using the default logger
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warnf logs a formatted warning message using the default logger
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Error logs an error message using the default logger
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorf logs a formatted error message using the default logger
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Fatalf logs a formatted fatal message using the default logger and exits
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}
