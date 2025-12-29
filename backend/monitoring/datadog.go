package monitoring

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Tracker provides custom metrics and tracing for business events
type Tracker struct {
	enabled     bool
	serviceName string
}

// NewTracker creates a new monitoring tracker
func NewTracker(enabled bool, serviceName string) *Tracker {
	return &Tracker{
		enabled:     enabled,
		serviceName: serviceName,
	}
}

// IsEnabled returns whether Datadog tracking is enabled
func (t *Tracker) IsEnabled() bool {
	return t.enabled
}

// TrackUserRegistration tracks user registration events
func (t *Tracker) TrackUserRegistration(ctx context.Context, userID uint, email string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "user.registration",
		tracer.Tag("user.id", userID),
		tracer.Tag("user.email", email),
		tracer.Tag("event.type", "user_registration"),
	)
	defer span.Finish()
}

// TrackTransactionCreated tracks when a transaction is created
func (t *Tracker) TrackTransactionCreated(ctx context.Context, userID uint, txType string, amount float64, category string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "transaction.created",
		tracer.Tag("user.id", userID),
		tracer.Tag("transaction.type", txType),
		tracer.Tag("transaction.amount", amount),
		tracer.Tag("transaction.category", category),
		tracer.Tag("event.type", "transaction_created"),
	)
	defer span.Finish()
}

// TrackAccountCreated tracks when an account is created
func (t *Tracker) TrackAccountCreated(ctx context.Context, userID uint, accountType string, accountName string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "account.created",
		tracer.Tag("user.id", userID),
		tracer.Tag("account.type", accountType),
		tracer.Tag("account.name", accountName),
		tracer.Tag("event.type", "account_created"),
	)
	defer span.Finish()
}

// TrackGoalContribution tracks goal contributions
func (t *Tracker) TrackGoalContribution(ctx context.Context, userID uint, goalID uint, amount float64) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "goal.contribution",
		tracer.Tag("user.id", userID),
		tracer.Tag("goal.id", goalID),
		tracer.Tag("contribution.amount", amount),
		tracer.Tag("event.type", "goal_contribution"),
	)
	defer span.Finish()
}

// TrackBudgetCreated tracks budget creation
func (t *Tracker) TrackBudgetCreated(ctx context.Context, userID uint, category string, amount float64, period string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "budget.created",
		tracer.Tag("user.id", userID),
		tracer.Tag("budget.category", category),
		tracer.Tag("budget.amount", amount),
		tracer.Tag("budget.period", period),
		tracer.Tag("event.type", "budget_created"),
	)
	defer span.Finish()
}

// TrackCreditCardPayment tracks credit card payments
func (t *Tracker) TrackCreditCardPayment(ctx context.Context, userID uint, cardID uint, amount float64) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "credit_card.payment",
		tracer.Tag("user.id", userID),
		tracer.Tag("card.id", cardID),
		tracer.Tag("payment.amount", amount),
		tracer.Tag("event.type", "credit_card_payment"),
	)
	defer span.Finish()
}

// TrackDebtPayment tracks debt payments
func (t *Tracker) TrackDebtPayment(ctx context.Context, userID uint, debtID uint, amount float64) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "debt.payment",
		tracer.Tag("user.id", userID),
		tracer.Tag("debt.id", debtID),
		tracer.Tag("payment.amount", amount),
		tracer.Tag("event.type", "debt_payment"),
	)
	defer span.Finish()
}

// TrackReconciliation tracks account reconciliations
func (t *Tracker) TrackReconciliation(ctx context.Context, userID uint, accountID uint, difference float64) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "account.reconciliation",
		tracer.Tag("user.id", userID),
		tracer.Tag("account.id", accountID),
		tracer.Tag("reconciliation.difference", difference),
		tracer.Tag("event.type", "account_reconciliation"),
	)
	defer span.Finish()
}

// TrackError tracks application errors with context
func (t *Tracker) TrackError(ctx context.Context, err error, operation string, metadata map[string]interface{}) {
	if !t.enabled || err == nil {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "error.occurred",
		tracer.Tag("error.message", err.Error()),
		tracer.Tag("error.operation", operation),
		tracer.Tag("event.type", "error"),
	)

	// Add metadata as tags
	for key, value := range metadata {
		span.SetTag(fmt.Sprintf("error.%s", key), value)
	}

	span.SetTag("error", err)
	defer span.Finish()
}

// MeasureOperation measures the duration of an operation
func (t *Tracker) MeasureOperation(ctx context.Context, operationName string, tags map[string]interface{}) func() {
	if !t.enabled {
		return func() {}
	}

	start := time.Now()
	span, _ := tracer.StartSpanFromContext(ctx, operationName)

	// Add custom tags
	for key, value := range tags {
		span.SetTag(key, value)
	}

	return func() {
		duration := time.Since(start)
		span.SetTag("duration.ms", duration.Milliseconds())
		span.Finish()
	}
}

// TrackReportGeneration tracks report generation events
func (t *Tracker) TrackReportGeneration(ctx context.Context, userID uint, reportType string, dateRange string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "report.generated",
		tracer.Tag("user.id", userID),
		tracer.Tag("report.type", reportType),
		tracer.Tag("report.date_range", dateRange),
		tracer.Tag("event.type", "report_generated"),
	)
	defer span.Finish()
}

// TrackBulkImport tracks bulk transaction imports
func (t *Tracker) TrackBulkImport(ctx context.Context, userID uint, count int, source string) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "transaction.bulk_import",
		tracer.Tag("user.id", userID),
		tracer.Tag("import.count", count),
		tracer.Tag("import.source", source),
		tracer.Tag("event.type", "bulk_import"),
	)
	defer span.Finish()
}

// TrackAPICall tracks external API calls
func (t *Tracker) TrackAPICall(ctx context.Context, apiName string, endpoint string, duration time.Duration, statusCode int) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "external.api_call",
		tracer.Tag("api.name", apiName),
		tracer.Tag("api.endpoint", endpoint),
		tracer.Tag("api.duration_ms", duration.Milliseconds()),
		tracer.Tag("api.status_code", statusCode),
		tracer.Tag("event.type", "api_call"),
	)
	defer span.Finish()
}

// TrackCacheOperation tracks cache hits and misses
func (t *Tracker) TrackCacheOperation(ctx context.Context, operation string, key string, hit bool) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "cache.operation",
		tracer.Tag("cache.operation", operation),
		tracer.Tag("cache.key", key),
		tracer.Tag("cache.hit", hit),
		tracer.Tag("event.type", "cache_operation"),
	)
	defer span.Finish()
}

// TrackDatabaseQuery tracks custom database query metrics
func (t *Tracker) TrackDatabaseQuery(ctx context.Context, queryType string, table string, duration time.Duration, rowsAffected int64) {
	if !t.enabled {
		return
	}

	span, _ := tracer.StartSpanFromContext(ctx, "database.custom_query",
		tracer.Tag("db.query_type", queryType),
		tracer.Tag("db.table", table),
		tracer.Tag("db.duration_ms", duration.Milliseconds()),
		tracer.Tag("db.rows_affected", rowsAffected),
		tracer.Tag("event.type", "database_query"),
	)
	defer span.Finish()
}

// AddSpanTags adds custom tags to the current span
func (t *Tracker) AddSpanTags(ctx context.Context, tags map[string]interface{}) {
	if !t.enabled {
		return
	}

	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return
	}

	for key, value := range tags {
		span.SetTag(key, value)
	}
}

// SetUser sets user information on the current span
func (t *Tracker) SetUser(ctx context.Context, userID uint, email string) {
	if !t.enabled {
		return
	}

	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return
	}

	span.SetTag("user.id", userID)
	span.SetTag("user.email", email)
}
