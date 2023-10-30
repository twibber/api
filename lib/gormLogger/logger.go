package gormLogger

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func New() *Logger {
	return &Logger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *Logger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Infof(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Warnf(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...any) {
	log.WithContext(ctx).Errorf(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		log.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if l.Debug {
		log.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}
