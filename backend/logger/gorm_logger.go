package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses our custom logger
type GormLogger struct {
	logger                    *Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger
func NewGormLogger(logger *Logger, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logger:                    logger,
		logLevel:                  gormlogger.Info,
		slowThreshold:             slowThreshold,
		ignoreRecordNotFoundError: true,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Infof(ctx, msg, data...)
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Warnf(ctx, msg, data...)
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Errorf(ctx, msg, data...)
	}
}

// Trace logs SQL queries with execution time
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		// Log error with query details
		l.logger.Errorf(ctx, "SQL Error | duration: %s | rows: %d | sql: %s | error: %v",
			elapsed, rows, sql, err)

	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		// Log slow query
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		l.logger.Warnf(ctx, "%s | duration: %s | rows: %d | sql: %s",
			slowLog, elapsed, rows, sql)

	case l.logLevel >= gormlogger.Info:
		// Log normal query
		l.logger.Infof(ctx, "SQL Query | duration: %s | rows: %d | sql: %s",
			elapsed, rows, sql)
	}
}

// SetLogLevel sets the GORM log level
func (l *GormLogger) SetLogLevel(level gormlogger.LogLevel) {
	l.logLevel = level
}

// SetSlowThreshold sets the slow query threshold
func (l *GormLogger) SetSlowThreshold(threshold time.Duration) {
	l.slowThreshold = threshold
}
