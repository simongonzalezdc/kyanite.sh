package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Kyanite/noise/internal/config"
	errutil "github.com/Kyanite/noise/internal/errutil"
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

// Logger wraps the standard log package with additional functionality
type Logger struct {
	*log.Logger
	level      LogLevel
	showCaller bool
	callerSkip int
	tuiMode    bool    // When true, suppress all output to avoid corrupting TUI
	logFile    *os.File // File handle for log file (if any)
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
	var output io.Writer = cfg.Output

	// If log file is specified, create or append to it
	var logFile *os.File
	if cfg.LogFile != "" {
		logDir := filepath.Dir(cfg.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, errutil.Wrap(err, "create log directory")
		}

		var err error
		logFile, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, errutil.Wrap(err, "open log file")
		}

		// Use multi-writer to write to both file and stdout
		output = io.MultiWriter(cfg.Output, logFile)
	}

	logger := &Logger{
		Logger:     log.New(output, "", log.LstdFlags|log.Lmicroseconds),
		level:      cfg.Level,
		showCaller: cfg.ShowCaller,
		callerSkip: 2,
		logFile:    logFile,
	}

	return logger, nil
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
}

// SetShowCaller sets whether to show caller information
func (l *Logger) SetShowCaller(show bool) {
	l.showCaller = show
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
	if l.level <= DEBUG {
		l.log(DEBUG, v...)
	}
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logf(DEBUG, format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, v...)
	}
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logf(INFO, format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, v...)
	}
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logf(WARN, format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, v...)
	}
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logf(ERROR, format, v...)
	}
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(v ...interface{}) {
	if l.level <= FATAL {
		l.log(FATAL, v...)
	}
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits the program
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level <= FATAL {
		l.logf(FATAL, format, v...)
	}
	os.Exit(1)
}

// log logs a message with the specified level
func (l *Logger) log(level LogLevel, v ...interface{}) {
	// Suppress output in TUI mode to avoid corrupting terminal
	if l.tuiMode {
		return
	}

	prefix := fmt.Sprintf("[%s] ", level.String())

	var message string
	if l.showCaller {
		file, line := l.getCaller()
		message = fmt.Sprintf("%s%s:%d ", prefix, file, line)
	} else {
		message = prefix
	}

	message += fmt.Sprint(v...)
	l.Logger.Print(message)
}

// logf logs a formatted message with the specified level
func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	// Suppress output in TUI mode to avoid corrupting terminal
	if l.tuiMode {
		return
	}

	prefix := fmt.Sprintf("[%s] ", level.String())

	var message string
	if l.showCaller {
		file, line := l.getCaller()
		message = fmt.Sprintf("%s%s:%d ", prefix, file, line)
	} else {
		message = prefix
	}

	message += fmt.Sprintf(format, v...)
	if l.Logger != nil {
		l.Logger.Print(message)
	}
}

// getCaller returns the file and line number of the caller
func (l *Logger) getCaller() (string, int) {
	_, file, line, ok := runtime.Caller(l.callerSkip)
	if !ok {
		return "unknown", 0
	}

	// Extract just the filename (not the full path)
	parts := strings.Split(file, string(os.PathSeparator))
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}

	return file, line
}

// Global logger instance
var defaultLogger *Logger

// init initializes the default logger
func init() {
	cfg := DefaultConfig()
	logger, err := New(cfg)
	if err != nil {
		// Fallback to standard logger if initialization fails
		// CRITICAL: Use stderr to avoid corrupting Bubble Tea's stdout rendering
		defaultLogger = &Logger{
			Logger: log.New(os.Stderr, "", log.LstdFlags),
			level:  INFO,
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
