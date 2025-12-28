package logger

import (
	"io"
	"os"
	"sync/atomic"
)

// loggerInstance manages logger configuration with atomic operations
type loggerInstance struct {
	writer      atomic.Pointer[io.Writer]
	logFile     atomic.Pointer[os.File]
	environment atomic.Pointer[string]
	logToFile   atomic.Bool
	initialized atomic.Bool
}

// init initializes the logger instance
func (l *loggerInstance) init(isLogToFile bool, filePath, env string) error {
	l.logToFile.Store(isLogToFile)

	// Store environment atomically
	envCopy := env
	l.environment.Store(&envCopy)

	var writer io.Writer

	if isLogToFile {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}

		l.logFile.Store(file)
		writer = file
	} else {
		writer = os.Stdout
	}

	l.writer.Store(&writer)
	l.initialized.Store(true)

	return nil
}

// close closes the logger and releases resources
func (l *loggerInstance) close() error {
	if !l.initialized.Load() {
		return nil
	}

	file := l.logFile.Load()
	if file != nil {
		return (*file).Close()
	}
	return nil
}

// isInitialized checks if logger is initialized
func (l *loggerInstance) isInitialized() bool {
	return l.initialized.Load()
}

// getWriter returns the current writer
func (l *loggerInstance) getWriter() io.Writer {
	writer := l.writer.Load()
	if writer == nil {
		return nil
	}
	return *writer
}

// isDev checks if running in development mode
func (l *loggerInstance) isDev() bool {
	env := l.environment.Load()
	if env == nil {
		return false
	}
	return *env == "dev"
}
