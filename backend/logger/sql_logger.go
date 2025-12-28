package logger

import (
	"context"
	"time"

	sqlLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// CustomLogger implements gorm's logger.Interface
type CustomLogger struct {
	Config sqlLogger.Config
}

// NewCustomLogger creates a new custom GORM logger
func NewCustomLogger(config sqlLogger.Config) *CustomLogger {
	return &CustomLogger{
		Config: config,
	}
}

// NewDefaultCustomLogger creates a custom logger with default config
func NewDefaultCustomLogger() *CustomLogger {
	return &CustomLogger{
		Config: sqlLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  sqlLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	}
}

// LogMode sets the log level
func (l *CustomLogger) LogMode(level sqlLogger.LogLevel) sqlLogger.Interface {
	newLogger := *l
	newLogger.Config.LogLevel = level
	return &newLogger
}

// Info logs info level messages
func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= sqlLogger.Info {
		Infof(ctx, msg, data...)
	}
}

// Warn logs warning level messages
func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= sqlLogger.Warn {
		Warnf(ctx, msg, data...)
	}
}

// Error logs error level messages
func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= sqlLogger.Error {
		Errorf(ctx, msg, data...)
	}
}

// Trace logs SQL queries with execution time and error handling
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Config.LogLevel <= sqlLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1e6

	sql, rows := fc()
	fileWithLine := utils.FileWithLineNum()

	switch {
	case err != nil && l.Config.LogLevel >= sqlLogger.Error:
		if !l.Config.IgnoreRecordNotFoundError || err.Error() != "record not found" {
			Errorf(ctx, "[SQL ERROR] time=%.3fms file=%s sql=%s rows=%d error=%v",
				elapsedMs, fileWithLine, sql, rows, err)
		}

	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Config.LogLevel >= sqlLogger.Warn:
		Warnf(ctx, "[SLOW SQL] time=%.3fms file=%s sql=%s rows=%d threshold=%v",
			elapsedMs, fileWithLine, sql, rows, l.Config.SlowThreshold)

	case l.Config.LogLevel == sqlLogger.Info:
		Infof(ctx, "[SQL] time=%.3fms file=%s sql=%s rows=%d",
			elapsedMs, fileWithLine, sql, rows)
	}
}
