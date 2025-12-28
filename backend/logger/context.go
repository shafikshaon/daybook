package logger

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContextToGinContext converts a context.Context to *gin.Context
func ContextToGinContext(ctx context.Context) *gin.Context {
	// Create a dummy request
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(ctx)

	ginCtx := &gin.Context{
		Request: req,
	}

	if data, ok := extractLogData(ctx); ok {
		if ginCtx.Keys == nil {
			ginCtx.Keys = make(map[any]any)
		}
		ginCtx.Set("traceID", data.TraceID)
	}

	return ginCtx
}

// GinToContext converts *gin.Context to context.Context
func GinToContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// extractLogData extracts log data from context
func extractLogData(ctx context.Context) (logData, bool) {
	data, ok := ctx.Value(logTracingKey).(logData)
	return data, ok
}
