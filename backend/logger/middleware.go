package logger

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// generateID generates a new UUID v7 without dashes
func generateID() string {
	ID, err := uuid.NewV7()
	if err != nil {
		Errorf(context.Background(), "cannot generate uuid: %v", err)
		return ""
	}

	return strings.ReplaceAll(ID.String(), "-", "")
}

// ModifyContext adds trace ID to gin context
func ModifyContext(c *gin.Context) {
	traceID := c.GetHeader("X-Trace-ID")
	if traceID == "" {
		traceID = c.GetHeader("X-Trace-Id")
	}
	if traceID == "" {
		traceID = generateID()
	}

	data := logData{
		TraceID: traceID,
	}

	ctx := context.WithValue(c.Request.Context(), logTracingKey, data)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// ModifyGRPCContext adds trace ID to gRPC context
func ModifyGRPCContext(ctx context.Context) context.Context {
	var traceID string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("x-trace-id")
		if len(values) > 0 {
			traceID = values[0]
		}
	}

	if traceID == "" {
		traceID = generateID()
	}

	data := logData{
		TraceID: traceID,
	}

	return context.WithValue(ctx, logTracingKey, data)
}

// CreateContext creates a new context with the given trace ID
func CreateContext(traceID string) context.Context {
	if traceID == "" {
		traceID = generateID()
	}

	data := logData{
		TraceID: traceID,
	}

	return context.WithValue(context.Background(), logTracingKey, data)
}
