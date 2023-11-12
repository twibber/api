package gormLogger

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus" // Logrus for structured logging
	"gorm.io/gorm"                   // GORM for ORM functionality
	gormLogger "gorm.io/gorm/logger" // Logger interface for GORM
	"gorm.io/gorm/utils"             // Utilities for GORM
)

// Logger struct to customise GORM logging behaviour.
type Logger struct {
	SlowThreshold         time.Duration // Duration after which a query is considered slow
	SourceField           string        // Field for the source of the log
	SkipErrRecordNotFound bool          // Flag to skip logging "record not found" errors
	Debug                 bool          // Flag to enable debug logging
}

// New creates a new instance of Logger with default settings.
func New() *Logger {
	return &Logger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

// LogMode for GORM logger interface, currently a no-op.
func (l *Logger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info logs information level messages.
func (l *Logger) Info(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Infof(s, args...)
}

// Warn logs warning level messages.
func (l *Logger) Warn(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Warnf(s, args...)
}

// Error logs error level messages.
func (l *Logger) Error(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Errorf(s, args...)
}

// Trace logs SQL queries, their execution time, and errors if any.
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin) // Calculates elapsed time for the query
	sql, _ := fc()               // Retrieves the SQL query and affected rows
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum() // Adds source field if specified
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		log.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed) // Logs errors except skipped "record not found"
		return
	}

	// Warns if query execution time exceeds the slow threshold
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	// Debug logs for every query if debug is enabled
	if l.Debug {
		log.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}
