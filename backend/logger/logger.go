package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

// Context keys for trace information
type contextKey string

const (
	TraceIDKey contextKey = "trace_id"
	SpanIDKey  contextKey = "span_id"
	UserIDKey  contextKey = "user_id"
)

// Logger represents the custom logger
type Logger struct {
	level      LogLevel
	output     io.Writer
	useColors  bool
	timeFormat string
}

var (
	// DefaultLogger is the default logger instance
	DefaultLogger *Logger
)

func init() {
	DefaultLogger = New(InfoLevel, os.Stdout, true)
}

// New creates a new Logger instance
func New(level LogLevel, output io.Writer, useColors bool) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		useColors:  useColors,
		timeFormat: "2006-01-02 15:04:05.000",
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// getCaller returns caller information (file, line, function)
func getCaller(skip int) (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}

	// Get short file path (relative to project root)
	file = filepath.Base(file)

	// Get function name
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		function = fn.Name()
		// Extract just the function name without package path
		parts := strings.Split(function, "/")
		if len(parts) > 0 {
			function = parts[len(parts)-1]
		}
	} else {
		function = "unknown"
	}

	return
}

// formatLog formats the log message with all context information
func (l *Logger) formatLog(ctx context.Context, level LogLevel, format string, args ...interface{}) string {
	var builder strings.Builder

	// Get caller information (skip 3 levels: formatLog -> log method -> user code)
	file, line, function := getCaller(3)

	// Get timestamp
	timestamp := time.Now().Format(l.timeFormat)

	// Get level color
	levelColor := l.getLevelColor(level)
	levelStr := level.String()

	// Extract context information
	traceID := extractContextValue(ctx, TraceIDKey, "-")
	spanID := extractContextValue(ctx, SpanIDKey, "-")
	userID := extractContextValue(ctx, UserIDKey, "-")

	// Format message
	message := fmt.Sprintf(format, args...)

	// Build log string
	if l.useColors {
		builder.WriteString(fmt.Sprintf("%s[%s]%s ", colorGray, timestamp, colorReset))
		builder.WriteString(fmt.Sprintf("%s[%-5s]%s ", levelColor, levelStr, colorReset))
		builder.WriteString(fmt.Sprintf("%s[trace:%s]%s ", colorCyan, traceID, colorReset))
		builder.WriteString(fmt.Sprintf("%s[span:%s]%s ", colorCyan, spanID, colorReset))
		builder.WriteString(fmt.Sprintf("%s[user:%s]%s ", colorPurple, userID, colorReset))
		builder.WriteString(fmt.Sprintf("%s[%s:%d]%s ", colorBlue, file, line, colorReset))
		builder.WriteString(fmt.Sprintf("%s[%s]%s ", colorGreen, function, colorReset))
		builder.WriteString(fmt.Sprintf("- %s", message))
	} else {
		builder.WriteString(fmt.Sprintf("[%s] [%-5s] [trace:%s] [span:%s] [user:%s] [%s:%d] [%s] - %s",
			timestamp, levelStr, traceID, spanID, userID, file, line, function, message))
	}

	return builder.String()
}

// getLevelColor returns the color for a log level
func (l *Logger) getLevelColor(level LogLevel) string {
	if !l.useColors {
		return ""
	}
	switch level {
	case DebugLevel:
		return colorGray
	case InfoLevel:
		return colorGreen
	case WarnLevel:
		return colorYellow
	case ErrorLevel:
		return colorRed
	default:
		return colorReset
	}
}

// extractContextValue extracts a value from context or returns default
func extractContextValue(ctx context.Context, key contextKey, defaultValue string) string {
	if ctx == nil {
		return defaultValue
	}
	if value := ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// log is the internal logging method
func (l *Logger) log(ctx context.Context, level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	logMessage := l.formatLog(ctx, level, format, args...)
	log.Println(logMessage)
}

// Debugf logs a debug message
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, args...)
}

// Infof logs an info message
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, InfoLevel, format, args...)
}

// Warnf logs a warning message
func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WarnLevel, format, args...)
}

// Errorf logs an error message
func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ErrorLevel, format, args...)
}

// Package-level convenience functions using DefaultLogger

// Debugf logs a debug message using the default logger
func Debugf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Debugf(ctx, format, args...)
}

// Infof logs an info message using the default logger
func Infof(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Infof(ctx, format, args...)
}

// Warnf logs a warning message using the default logger
func Warnf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Warnf(ctx, format, args...)
}

// Errorf logs an error message using the default logger
func Errorf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Errorf(ctx, format, args...)
}

// SetLevel sets the log level for the default logger
func SetLevel(level LogLevel) {
	DefaultLogger.SetLevel(level)
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithSpanID adds a span ID to the context
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, SpanIDKey, spanID)
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetTraceID retrieves the trace ID from context
func GetTraceID(ctx context.Context) string {
	return extractContextValue(ctx, TraceIDKey, "")
}

// GetSpanID retrieves the span ID from context
func GetSpanID(ctx context.Context) string {
	return extractContextValue(ctx, SpanIDKey, "")
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) string {
	return extractContextValue(ctx, UserIDKey, "")
}
