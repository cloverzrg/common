package logger

import (
	"context"
	"errors"
	"fmt"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	LogLevel                  logger2.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	Colorful                  bool
}

func GetGormLogger() *GormLogger {
	return &GormLogger{
		LogLevel:                  logger2.Warn,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             15 * time.Millisecond,
		Colorful:                  false,
	}
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger2.LogLevel) logger2.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Info {
		logger.Infof(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Warn {
		logger.Warnf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Error {
		logger.Errorf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger2.Silent {
		return
	}
	traceWarnStr := "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr := "%s %s\n[%.3fms] [rows:%v] %s"
	traceStr := "%s\n[%.3fms] [rows:%v] %s"

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger2.Error && (!errors.Is(err, logger2.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger2.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger2.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Tracef(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Tracef(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
