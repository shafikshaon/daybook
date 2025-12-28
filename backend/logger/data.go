package logger

import "context"

type contextKey string

const (
	logTracingKey contextKey = "LogTracingKey"
)

// logData holds tracing information
type logData struct {
	TraceID string
}

// GetTraceID returns the trace ID associated with the context
func GetTraceID(ctx context.Context) string {
	if data, ok := ctx.Value(logTracingKey).(logData); ok {
		return data.TraceID
	}
	return ""
}
