# Comprehensive Datadog Monitoring Guide

This document explains what gets tracked when Datadog APM is enabled in the Daybook backend.

## Overview

When `DD_ENABLED=true`, the application tracks **everything** across all layers:
- HTTP requests and responses
- Database queries (PostgreSQL via GORM)
- Cache operations (Redis)
- Business events and user activities
- Errors and exceptions
- Custom performance metrics

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Your Application Code                       │
├─────────────────────────────────────────────────────────┤
│  HTTP Layer (Gin)                                       │
│    → All API endpoints automatically traced            │
├─────────────────────────────────────────────────────────┤
│  Database Layer (GORM)                                  │
│    → All SQL queries automatically traced              │
├─────────────────────────────────────────────────────────┤
│  Cache Layer (Redis)                                    │
│    → All Redis commands automatically traced           │
├─────────────────────────────────────────────────────────┤
│  Custom Business Metrics (monitoring.Tracker)           │
│    → Manual instrumentation for business events        │
└─────────────────────────────────────────────────────────┘
                        ↓
              Datadog Agent (localhost:8126)
                        ↓
              Datadog Cloud Platform
```

## 1. HTTP Request Tracking (Automatic)

**What**: Every HTTP request to your API
**How**: Gin middleware (`gintrace.Middleware`)
**Location**: `main.go:98`

### Tracked Automatically:
- **Endpoint**: Route pattern (e.g., `/api/v1/transactions/:id`)
- **Method**: HTTP method (GET, POST, PUT, DELETE, etc.)
- **Status Code**: Response status (200, 404, 500, etc.)
- **Duration**: Total request processing time
- **URL**: Full request URL
- **Headers**: Selected headers (User-Agent, etc.)
- **Query Parameters**: URL query parameters
- **Request Size**: Size of request body
- **Response Size**: Size of response body

### Example Traces:
```
POST /api/v1/auth/login               → 201 Created (145ms)
GET  /api/v1/transactions             → 200 OK (87ms)
POST /api/v1/transactions             → 201 Created (234ms)
PUT  /api/v1/accounts/123             → 200 OK (156ms)
DELETE /api/v1/transactions/456       → 204 No Content (98ms)
```

## 2. Database Query Tracking (Automatic)

**What**: All PostgreSQL database operations
**How**: GORM plugin (`gormtrace.NewTracePlugin`)
**Location**: `database/database.go:45`

### Tracked Automatically:
- **SQL Query**: Full SQL statement (sanitized)
- **Table**: Database table being queried
- **Operation Type**: SELECT, INSERT, UPDATE, DELETE
- **Duration**: Query execution time
- **Rows Affected**: Number of rows returned/modified
- **Database Name**: Database being queried
- **Connection Pool**: Connection pool metrics
- **Errors**: SQL errors and constraint violations

### Example Traces:
```sql
SELECT * FROM transactions WHERE user_id = ? AND date >= ? AND date <= ?
→ Duration: 23ms, Rows: 145

INSERT INTO accounts (user_id, name, type, balance) VALUES (?, ?, ?, ?)
→ Duration: 12ms, Rows: 1

UPDATE transactions SET amount = ?, description = ? WHERE id = ?
→ Duration: 8ms, Rows: 1

DELETE FROM transactions WHERE id = ? AND user_id = ?
→ Duration: 6ms, Rows: 1
```

### Performance Insights:
- Identify slow queries (N+1 problems, missing indexes)
- Track query frequency
- Monitor connection pool saturation
- Detect deadlocks and lock contention

## 3. Redis Cache Tracking (Automatic)

**What**: All Redis operations
**How**: Redis wrapper (`redistrace.NewClient`)
**Location**: `database/database.go:103`

### Tracked Automatically:
- **Command**: Redis command (GET, SET, DEL, etc.)
- **Key**: Cache key being accessed
- **Duration**: Command execution time
- **Result**: Success/failure
- **Connection**: Connection pool metrics

### Example Traces:
```
GET user:sessions:123           → Duration: 2ms (HIT)
SET user:profile:456 EX 3600    → Duration: 3ms (SUCCESS)
DEL transaction:cache:*         → Duration: 5ms (10 keys deleted)
HGETALL account:balances:789    → Duration: 4ms (HIT)
```

### Cache Metrics:
- Cache hit ratio
- Cache miss rate
- Average lookup time
- Eviction rates

## 4. Custom Business Event Tracking (Manual)

**What**: Important business events and user actions
**How**: `monitoring.Tracker` utility
**Location**: `monitoring/datadog.go`

The `monitoring.Tracker` is available in the container and can be used in services:

```go
// Example usage in a service
func (s *transactionService) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
    // Create the transaction
    err := s.repo.Create(ctx, tx)
    if err != nil {
        return err
    }

    // Track the business event
    s.monitor.TrackTransactionCreated(ctx, tx.UserID, tx.Type, tx.Amount, tx.Category)

    return nil
}
```

### Available Tracking Methods:

#### User Events
```go
// User registration
tracker.TrackUserRegistration(ctx, userID, email)
```

#### Transaction Events
```go
// Transaction created
tracker.TrackTransactionCreated(ctx, userID, txType, amount, category)

// Bulk import
tracker.TrackBulkImport(ctx, userID, count, source)
```

#### Account Events
```go
// Account created
tracker.TrackAccountCreated(ctx, userID, accountType, accountName)

// Reconciliation performed
tracker.TrackReconciliation(ctx, userID, accountID, difference)
```

#### Financial Operations
```go
// Goal contribution
tracker.TrackGoalContribution(ctx, userID, goalID, amount)

// Budget created
tracker.TrackBudgetCreated(ctx, userID, category, amount, period)

// Credit card payment
tracker.TrackCreditCardPayment(ctx, userID, cardID, amount)

// Debt payment
tracker.TrackDebtPayment(ctx, userID, debtID, amount)
```

#### Reports & Analytics
```go
// Report generated
tracker.TrackReportGeneration(ctx, userID, reportType, dateRange)
```

#### Error Tracking
```go
// Track application errors
tracker.TrackError(ctx, err, "create_transaction", map[string]interface{}{
    "user_id": userID,
    "amount": amount,
    "account_id": accountID,
})
```

#### Performance Measurement
```go
// Measure operation duration
done := tracker.MeasureOperation(ctx, "generate_monthly_report", map[string]interface{}{
    "user_id": userID,
    "month": "2024-01",
})
defer done()

// ... expensive operation ...
```

#### Custom Span Tags
```go
// Add custom tags to current span
tracker.AddSpanTags(ctx, map[string]interface{}{
    "feature_flag": "new_ui_enabled",
    "experiment_id": "exp_123",
    "user_segment": "premium",
})

// Set user context
tracker.SetUser(ctx, userID, email)
```

## 5. Error Tracking

All errors are automatically captured with:
- **Error Message**: Full error description
- **Stack Trace**: Where the error occurred
- **Context**: Request/user/session information
- **Tags**: Custom tags for filtering

## 6. Service Map & Dependencies

Datadog automatically builds a service map showing:
- **daybook-backend**: Main API service
- **postgres**: Database dependency
- **daybook-backend-redis**: Redis cache dependency
- **External APIs**: Any third-party API calls

## 7. Key Performance Indicators (KPIs)

### Application Metrics
- **Request Rate**: Requests per second
- **Error Rate**: % of failed requests
- **Latency**: P50, P75, P95, P99 response times
- **Throughput**: Data processed per second

### Database Metrics
- **Query Rate**: Queries per second
- **Slow Queries**: Queries > 1000ms
- **Connection Pool**: Active/idle connections
- **Lock Contention**: Database lock wait times

### Cache Metrics
- **Hit Rate**: % of successful cache lookups
- **Miss Rate**: % of cache misses
- **Latency**: Average Redis command time
- **Memory Usage**: Redis memory consumption

### Business Metrics
- **User Registrations**: New users per hour/day
- **Transactions Created**: Financial transactions per hour
- **Active Users**: Unique users making requests
- **Feature Usage**: Which endpoints are most used

## 8. Custom Dashboards

You can create custom dashboards in Datadog to visualize:

### Financial Overview Dashboard
```
┌──────────────────────────────────────────┐
│ Transactions Created (24h)    1,234      │
│ Total Transaction Amount      $45,678    │
│ Active Users                  567        │
│ API Request Rate             120 req/s   │
└──────────────────────────────────────────┘
```

### Performance Dashboard
```
┌──────────────────────────────────────────┐
│ P95 Response Time            145ms       │
│ Error Rate                   0.2%        │
│ Database Query Time          23ms        │
│ Cache Hit Rate               87%         │
└──────────────────────────────────────────┘
```

## 9. Alerts & Monitoring

Set up monitors for:

### Critical Alerts
- Error rate > 1% for 5 minutes
- P95 latency > 1000ms
- Database connections > 90% of pool
- Service downtime

### Warning Alerts
- Cache hit rate < 70%
- Slow queries > 500ms
- API rate limiting triggered
- Memory usage > 80%

## 10. Querying Traces

### In Datadog UI

**Filter by Service:**
```
service:daybook-backend
```

**Filter by User:**
```
user.id:123
```

**Find Slow Requests:**
```
service:daybook-backend duration:>1s
```

**Find Errors:**
```
service:daybook-backend status:error
```

**Find Database Issues:**
```
resource_name:"SELECT * FROM transactions" duration:>500ms
```

**Find Specific Events:**
```
event.type:transaction_created user.id:123
```

## 11. Best Practices

### DO:
✅ Enable Datadog in production and staging
✅ Set up alerts for critical metrics
✅ Review slow queries weekly
✅ Monitor error rates daily
✅ Use custom tags for feature flags
✅ Track business metrics
✅ Set appropriate sampling rates for high traffic

### DON'T:
❌ Log sensitive data (passwords, API keys, tokens)
❌ Track PII without proper scrubbing
❌ Create too many custom metrics (cost implications)
❌ Ignore high error rates
❌ Leave default alert thresholds
❌ Track health check endpoints (noise)

## 12. Privacy & Security

### Sensitive Data Scrubbing

Datadog automatically scrubs common patterns:
- Passwords
- Credit card numbers
- API keys
- Authorization headers

### Additional Configuration

To scrub custom fields, configure in Datadog Agent:

```yaml
# /etc/datadog-agent/datadog.yaml
apm_config:
  obfuscation:
    elasticsearch:
      enabled: true
    mongodb:
      enabled: true
    sql_exec_plan:
      enabled: true
    sql_exec_plan_normalize:
      enabled: true

  # Custom scrubbing rules
  replace_tags:
    - name: "email"
      pattern: "([a-zA-Z0-9_.+-]+)@([a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+)"
      repl: "${1}@***"
```

## 13. Cost Optimization

### Sampling

For high-traffic endpoints, use sampling:

```go
// In main.go, configure sampling rate
if cfg.Datadog.Enabled {
    tracer.Start(
        tracer.WithService(cfg.Datadog.ServiceName),
        tracer.WithEnv(cfg.Datadog.Environment),
        tracer.WithSamplingRules([]tracer.SamplingRule{
            // Sample 100% of errors
            {Service: "daybook-backend", Rate: 1.0, Tags: map[string]string{"status": "error"}},
            // Sample 10% of health checks
            {Service: "daybook-backend", Name: "GET /health", Rate: 0.1},
            // Sample 50% of other requests
            {Service: "daybook-backend", Rate: 0.5},
        }),
    )
}
```

### Filter Noisy Endpoints

Exclude health checks and metrics endpoints:

```go
router.Use(gintrace.Middleware(
    cfg.Datadog.ServiceName,
    gintrace.WithIgnoreRequest(func(r *http.Request) bool {
        return r.URL.Path == "/health" || r.URL.Path == "/metrics"
    }),
))
```

## 14. Troubleshooting

### No Traces Appearing

1. **Check agent status:**
   ```bash
   sudo datadog-agent status
   ```

2. **Verify application logs:**
   ```
   Datadog APM tracer initialized ✓
   Datadog GORM tracing enabled ✓
   Datadog Redis tracing enabled ✓
   ```

3. **Check network:**
   ```bash
   telnet localhost 8126
   ```

### High Memory Usage

```yaml
# /etc/datadog-agent/datadog.yaml
apm_config:
  max_memory: 500000000  # 500MB
  max_cpu_percent: 50
```

### Traces Not Linking

Ensure context is properly propagated:

```go
// GOOD - context propagated
func (s *service) CreateTransaction(ctx context.Context, tx *Tx) error {
    return s.repo.Create(ctx, tx)  // ✓ ctx passed through
}

// BAD - context lost
func (s *service) CreateTransaction(ctx context.Context, tx *Tx) error {
    return s.repo.Create(context.Background(), tx)  // ✗ new context
}
```

## 15. Additional Resources

- **Datadog APM Documentation**: https://docs.datadoghq.com/tracing/
- **Go Tracer Guide**: https://docs.datadoghq.com/tracing/setup_overview/setup/go/
- **Service Map**: https://app.datadoghq.com/apm/map
- **Trace Explorer**: https://app.datadoghq.com/apm/traces
- **Monitors**: https://app.datadoghq.com/monitors/manage

## 16. Example Queries for Common Scenarios

### Find users experiencing errors
```
service:daybook-backend status:error @user.id:*
```

### Slow transaction creations
```
resource_name:"POST /api/v1/transactions" duration:>1s
```

### Database query performance
```
db.operation:SELECT db.table:transactions @duration:>100ms
```

### Cache efficiency
```
resource_name:"redis.get" @cache.hit:false
```

### High-value transactions
```
event.type:transaction_created @transaction.amount:>10000
```

---

**Happy Monitoring! 🚀**

For support, check the Datadog documentation or reach out to the development team.
