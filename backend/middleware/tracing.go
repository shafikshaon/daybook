package middleware

import (
	"context"
	"fmt"

	"daybook-backend/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TracingMiddleware adds trace ID and span ID to the request context
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or extract trace ID from header
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Generate span ID for this request
		spanID := generateSpanID()

		// Create context with trace information
		ctx := logger.WithTraceID(c.Request.Context(), traceID)
		ctx = logger.WithSpanID(ctx, spanID)

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Set trace ID in response header for client tracking
		c.Header("X-Trace-ID", traceID)
		c.Header("X-Span-ID", spanID)

		// Store in gin context for easy access
		c.Set("traceID", traceID)
		c.Set("spanID", spanID)

		// Log request entry
		logger.Infof(ctx, "Request started: %s %s | client_ip: %s | user_agent: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP(), c.Request.UserAgent())

		// Process request
		c.Next()

		// Log request completion
		logger.Infof(ctx, "Request completed: %s %s | status: %d",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

// generateSpanID generates a unique span ID
func generateSpanID() string {
	return fmt.Sprintf("span-%s", uuid.New().String()[:8])
}

// GetContext returns the context from gin.Context with trace information
func GetContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// GetContextWithUserID returns the context with user ID added
func GetContextWithUserID(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// Try to get user ID from context
	userID, err := GetUserID(c)
	if err == nil {
		ctx = logger.WithUserID(ctx, userID.String())
	}

	return ctx
}
