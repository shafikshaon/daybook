# Custom Logger Module

This package provides a comprehensive logging solution for the Daybook backend application with request tracing, context-aware logging, and database query logging capabilities.

## Features

- **JSON Structured Logging**: All logs are output in JSON format for easy parsing and analysis
- **Context-Aware Logging**: Automatically includes trace ID, span ID, and user ID from context
- **Caller Information**: Automatically captures file name, line number, and function name
- **Log Levels**: Support for Debug, Info, Warn, and Error levels
- **Database Query Logging**: Custom GORM logger that logs all database queries with execution time
- **Request Tracing**: Middleware for automatic request tracing with unique trace and span IDs

## Usage

### Basic Logging

```go
import (
    "daybook-backend/logger"
    "daybook-backend/middleware"
)

func MyHandler(c *gin.Context) {
    ctx := middleware.GetContext(c)

    logger.Infof(ctx, "Processing request for user: %s", userID)
    logger.Debugf(ctx, "Detailed debug information: %+v", data)
    logger.Warnf(ctx, "Warning: %s", warningMessage)
    logger.Errorf(ctx, "Error occurred: %v", err)
}
```

### With User ID Context

```go
func AuthenticatedHandler(c *gin.Context) {
    // This automatically includes user ID in logs
    ctx := middleware.GetContextWithUserID(c)

    logger.Infof(ctx, "User accessed protected resource")
}
```

### Database Operations with Context

```go
func GetUser(c *gin.Context) {
    ctx := middleware.GetContextWithUserID(c)

    var user models.User
    // Use WithContext to ensure queries are logged with trace information
    if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
        logger.Errorf(ctx, "Failed to fetch user: %v", err)
        return
    }
}
```

## Log Format

Logs are output in structured JSON format with the following fields:

```json
{
  "timestamp": "2025-11-08T12:34:56.789123456Z",
  "level": "INFO",
  "trace_id": "abc-123-def-456",
  "span_id": "span-12345678",
  "user_id": "user-uuid-here",
  "file": "auth_handler.go",
  "line": 169,
  "function": "handlers.Login",
  "message": "User login completed successfully for user: johndoe (ID: user-uuid)"
}
```

Fields:
- **timestamp**: ISO 8601 format with nanosecond precision (RFC3339Nano)
- **level**: Log level (DEBUG, INFO, WARN, ERROR)
- **trace_id**: Unique request trace identifier (omitted if not set)
- **span_id**: Unique span identifier for the request (omitted if not set)
- **user_id**: Authenticated user ID (omitted if not set)
- **file**: Source file name where log was called
- **line**: Line number in source file
- **function**: Function name that generated the log
- **message**: The log message

## Log Levels

- **DEBUG**: Detailed information for debugging (e.g., hashing passwords, creating settings)
- **INFO**: General informational messages (e.g., handler entry/exit, successful operations)
- **WARN**: Warning messages (e.g., authentication failures, validation errors)
- **ERROR**: Error messages (e.g., database failures, critical errors)

## Request Tracing

The tracing middleware automatically:
1. Generates or extracts a trace ID from the `X-Trace-ID` header
2. Creates a unique span ID for each request
3. Logs request entry and exit
4. Adds trace information to response headers

## Database Query Logging

The custom GORM logger automatically logs:
- All SQL queries with execution time
- Slow queries (threshold: 200ms) with WARNING level
- Query errors with full context
- Number of rows affected

Example:
```
[2025-11-08 12:34:56.789] [INFO] [trace:abc-123] [span:span-12ab] [user:-] [gorm_logger.go:85] [logger.(*GormLogger).Trace] - SQL Query | duration: 15ms | rows: 1 | sql: SELECT * FROM users WHERE id = '...'
```

## Configuration

Set log level in `main.go`:

```go
// For production
logger.SetLevel(logger.InfoLevel)

// For development
logger.SetLevel(logger.DebugLevel)
```

## Best Practices

1. **Always use context**: Pass the context from `middleware.GetContext(c)` or `middleware.GetContextWithUserID(c)` to log functions
2. **Log at entry points**: Add `logger.Infof(ctx, "HandlerName - Entry")` at the start of handlers
3. **Log critical operations**: Log before and after important database operations
4. **Use appropriate levels**:
   - DEBUG for detailed flow
   - INFO for successful operations
   - WARN for expected errors (e.g., validation failures)
   - ERROR for unexpected errors
5. **Include relevant data**: Add IDs, counts, and other contextual information to log messages
6. **Use WithContext for all DB calls**: Always use `database.DB.WithContext(ctx)` to ensure query logging includes trace information

## Example Handler Pattern

```go
func MyHandler(c *gin.Context) {
    ctx := middleware.GetContextWithUserID(c)
    logger.Infof(ctx, "MyHandler - Entry")

    userID, err := middleware.GetUserID(c)
    if err != nil {
        logger.Warnf(ctx, "Unauthorized access: %v", err)
        utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
        return
    }

    logger.Debugf(ctx, "Processing request for user: %s", userID)

    var data models.MyModel
    if err := database.DB.WithContext(ctx).First(&data, id).Error; err != nil {
        logger.Errorf(ctx, "Failed to fetch data: %v", err)
        utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch data")
        return
    }

    logger.Infof(ctx, "Successfully processed request for user: %s", userID)
    utilities.SuccessResponse(c, data, "Success")
}
```

## Thread Safety

The logger is thread-safe and can be used concurrently across multiple goroutines.

## Performance

- Minimal overhead for log formatting
- Efficient caller information extraction using runtime.Caller
- No blocking I/O operations in critical path
