package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger        *zap.Logger
	lvl           logger.LogLevel
	slowThreshold time.Duration
}

func NewGorm(l *zap.Logger) *GormLogger {
	l = l.With(zap.String("service", "gorm")).WithOptions(zap.AddCallerSkip(1))
	return &GormLogger{
		logger:        l,
		lvl:           logger.Warn,
		slowThreshold: 100 * time.Millisecond,
	}
}

func (l *GormLogger) LogMode(lvl logger.LogLevel) logger.Interface {
	return &GormLogger{
		logger:        l.logger,
		slowThreshold: l.slowThreshold,
		lvl:           lvl,
	}
}

func (l *GormLogger) Info(_ context.Context, str string, args ...interface{}) {
	if l.lvl >= logger.Info {
		return
	}
	l.logger.Sugar().Debugf(str, args...)
}

func (l *GormLogger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.lvl >= logger.Warn {
		return
	}
	l.logger.Sugar().Warnf(str, args...)
}

func (l *GormLogger) Error(_ context.Context, str string, args ...interface{}) {
	if l.lvl >= logger.Error {
		return
	}
	l.logger.Sugar().Errorf(str, args...)
}

func (l *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.lvl <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil:
		// Propagated errors are not logged so as not to appear multiple times and be confusing
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.lvl >= logger.Warn:
		sql, rows := fc()
		msg := fmt.Sprintf("Slow query > %v", l.slowThreshold)
		// While the message is "warning" that a query was running for a long period of time,
		// it is not clearly actionable and does not indicate that the server will fail in the near future.
		// This should be logged at the info level instead.
		l.logger.Info(msg, zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.lvl >= logger.Info:
		sql, rows := fc()
		l.logger.Debug("SQL query", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
