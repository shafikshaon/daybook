package logger

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

const (
	LogFlag   = "LOG"
	InfoFlag  = "INFO"
	ErrorFlag = "ERROR"
	DebugFlag = "DEBUG"
	WarnFlag  = "WARN"
	TraceFlag = "TRACE"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Level      string `json:"level"`
	Function   string `json:"function"`
	TraceID    string `json:"traceId"`
	Message    string `json:"message"`
	Stacktrace string `json:"stackTrace,omitempty"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Timestamp  string `json:"timestamp"`
}

var (
	globalLogger loggerInstance
	initOnce     sync.Once
	initErr      error
)

// InitLogger initializes the global logger instance (thread-safe, called once)
func InitLogger(isLogToFile bool, filePath, env string) error {
	initOnce.Do(func() {
		if err := globalLogger.init(isLogToFile, filePath, env); err != nil {
			initErr = err
		}
	})
	return initErr
}

// Close closes the logger and releases resources
func Close() error {
	return globalLogger.close()
}

// Infof logs an info level message
func Infof(ctx interface{}, message string, args ...interface{}) {
	if !globalLogger.isInitialized() {
		return
	}
	globalLogger.log(getContext(ctx), InfoFlag, ColorGreen, message, args...)
}

// Errorf logs an error level message with stack trace
func Errorf(ctx interface{}, message string, args ...interface{}) {
	if !globalLogger.isInitialized() {
		return
	}
	globalLogger.log(getContext(ctx), ErrorFlag, ColorRed, message, args...)
}

// Warnf logs a warning level message
func Warnf(ctx interface{}, message string, args ...interface{}) {
	if !globalLogger.isInitialized() {
		return
	}
	globalLogger.log(getContext(ctx), WarnFlag, ColorYellow, message, args...)
}

// Debugf logs a debug level message
func Debugf(ctx interface{}, message string, args ...interface{}) {
	if !globalLogger.isInitialized() {
		return
	}
	globalLogger.log(getContext(ctx), DebugFlag, ColorMagenta, message, args...)
}

// Tracef logs a trace level message
func Tracef(ctx interface{}, message string, args ...interface{}) {
	if !globalLogger.isInitialized() {
		return
	}
	globalLogger.log(getContext(ctx), TraceFlag, ColorCyan, message, args...)
}

// log is the core logging function with optimizations
func (l *loggerInstance) log(ctx context.Context, level, color, message string, args ...interface{}) {
	writer := l.getWriter()
	if writer == nil {
		return
	}

	// Get pooled resources
	buf := bufferPool.Get().(*[]byte)
	entry := entryPool.Get().(*LogEntry)
	timeBuf := timeBufferPool.Get().(*[]byte)

	defer func() {
		*buf = (*buf)[:0]
		*timeBuf = (*timeBuf)[:0]
		bufferPool.Put(buf)
		entryPool.Put(entry)
		timeBufferPool.Put(timeBuf)
	}()

	// Format timestamp to buffer (zero-alloc)
	*timeBuf = time.Now().AppendFormat(*timeBuf, time.RFC3339Nano)

	// Get trace ID
	traceID := GetTraceID(ctx)

	// Format message
	var logMessage string
	if len(args) > 0 {
		logMessage = formatMessage(message, args...)
	} else {
		logMessage = message
	}

	// Get caller info
	funcName, fileName, lineNo := getCallerInfo()

	// Populate entry
	entry.Level = level
	entry.TraceID = traceID
	entry.Message = logMessage
	entry.Function = funcName
	entry.File = fileName
	entry.Line = lineNo
	entry.Timestamp = bytesToString(*timeBuf)

	// Get stack trace only for errors
	if level == ErrorFlag {
		entry.Stacktrace = getStackTrace()
	} else {
		entry.Stacktrace = ""
	}

	// Encode to JSON manually (avoids reflection)
	*buf = appendJSON(*buf, entry)

	// Apply color if dev environment
	if l.isDev() {
		writeWithColor(writer, *buf, color)
	} else {
		*buf = append(*buf, '\n')
		writer.Write(*buf)
	}
}

// appendJSON manually encodes LogEntry to JSON (zero reflection)
func appendJSON(buf []byte, entry *LogEntry) []byte {
	buf = append(buf, '{')

	buf = append(buf, `"level":"`...)
	buf = append(buf, entry.Level...)
	buf = append(buf, `","function":"`...)
	buf = append(buf, jsonEscape(entry.Function)...)
	buf = append(buf, `","traceId":"`...)
	buf = append(buf, entry.TraceID...)
	buf = append(buf, `","message":"`...)
	buf = append(buf, jsonEscape(entry.Message)...)
	buf = append(buf, '"')

	if entry.Stacktrace != "" {
		buf = append(buf, `,"stackTrace":"`...)
		buf = append(buf, jsonEscape(entry.Stacktrace)...)
		buf = append(buf, '"')
	}

	buf = append(buf, `,"file":"`...)
	buf = append(buf, jsonEscape(entry.File)...)
	buf = append(buf, `","line":`...)
	buf = appendInt(buf, entry.Line)
	buf = append(buf, `,"timestamp":"`...)
	buf = append(buf, entry.Timestamp...)
	buf = append(buf, `"}`...)

	return buf
}

// jsonEscape escapes special JSON characters
func jsonEscape(s string) string {
	// Fast path: check if escaping is needed
	needsEscape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' || c == '\n' || c == '\r' || c == '\t' {
			needsEscape = true
			break
		}
	}

	if !needsEscape {
		return s
	}

	// Slow path: escape
	sb := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		stringBuilderPool.Put(sb)
	}()

	sb.Grow(len(s) + 20)

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			if c < 0x20 {
				sb.WriteString(fmt.Sprintf(`\u%04x`, c))
			} else {
				sb.WriteByte(c)
			}
		}
	}

	return sb.String()
}

// appendInt converts integer to string and appends to buffer
func appendInt(buf []byte, i int) []byte {
	if i == 0 {
		return append(buf, '0')
	}

	if i < 0 {
		buf = append(buf, '-')
		i = -i
	}

	var tmp [20]byte
	pos := len(tmp)

	for i > 0 {
		pos--
		tmp[pos] = byte('0' + i%10)
		i /= 10
	}

	return append(buf, tmp[pos:]...)
}

// formatMessage formats the log message with arguments
func formatMessage(format string, args ...interface{}) string {
	// Fast path: no format verbs
	if !strings.Contains(format, "%") {
		if len(args) == 1 {
			return fmt.Sprint(args[0])
		}

		sb := stringBuilderPool.Get().(*strings.Builder)
		defer func() {
			sb.Reset()
			stringBuilderPool.Put(sb)
		}()

		sb.Grow(len(format) + len(args)*16)
		sb.WriteString(format)
		sb.WriteByte(' ')

		for i, arg := range args {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(sb, arg)
		}
		return sb.String()
	}

	// Has format verbs - use Sprintf
	return fmt.Sprintf(format, args...)
}

// getCallerInfo retrieves caller function, file, and line number
func getCallerInfo() (string, string, int) {
	var pcs [1]uintptr
	n := runtime.Callers(5, pcs[:])
	if n == 0 {
		return "unknown", "unknown", 0
	}

	frames := runtime.CallersFrames(pcs[:])
	frame, _ := frames.Next()

	funcName := extractFunctionName(frame.Function)

	return funcName, frame.File, frame.Line
}

// extractFunctionName extracts the short function name
func extractFunctionName(fullName string) string {
	// Find last slash
	if idx := strings.LastIndexByte(fullName, '/'); idx >= 0 {
		fullName = fullName[idx+1:]
	}

	// Find last dot
	if idx := strings.LastIndexByte(fullName, '.'); idx >= 0 {
		return fullName[idx+1:]
	}

	return fullName
}

// getStackTrace generates a stack trace for error logs
func getStackTrace() string {
	pcs := stackBufPool.Get().(*[]uintptr)
	defer stackBufPool.Put(pcs)

	n := runtime.Callers(5, *pcs)
	if n == 0 {
		return ""
	}

	sb := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		stringBuilderPool.Put(sb)
	}()

	sb.Grow(n * 120)

	frames := runtime.CallersFrames((*pcs)[:n])
	for {
		frame, more := frames.Next()

		funcName := extractFunctionName(frame.Function)

		sb.WriteString(funcName)
		sb.WriteString("\n\t")
		sb.WriteString(frame.File)
		sb.WriteByte(':')
		fmt.Fprintf(sb, "%d", frame.Line)
		sb.WriteByte('\n')

		if !more {
			break
		}
	}

	return sb.String()
}

// writeWithColor writes colored output for dev environment
func writeWithColor(w io.Writer, content []byte, color string) {
	size := len(color) + len(content) + len(ColorReset) + 1

	buf := make([]byte, 0, size)
	buf = append(buf, color...)
	buf = append(buf, content...)
	buf = append(buf, ColorReset...)
	buf = append(buf, '\n')

	w.Write(buf)
}

// bytesToString converts byte slice to string without allocation
func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// getContext extracts context.Context from various types
func getContext(ctx interface{}) context.Context {
	switch v := ctx.(type) {
	case *gin.Context:
		return v.Request.Context()
	case context.Context:
		return v
	default:
		return context.Background()
	}
}
