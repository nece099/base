package dblogger

import (
	"context"
	"time"

	glogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel glogger.LogLevel
}

func (l *GormLogger) LogMode(glogger.LogLevel) glogger.Interface {
	return l
}

func (l GormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel == glogger.Info {
		Log.Infof(format, args...)
	}
}
func (l GormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel == glogger.Info {
		Log.Warnf(format, args...)
	}
}
func (l GormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel == glogger.Info {
		Log.Errorf(format, args...)
	}
}
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == glogger.Info {
		elapsed := time.Since(begin)
		sql, rows := fc()

		if err == nil {
			Log.Debugf("sql log : time[%.3f],rows[%v]\n%v", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		} else {
			Log.Debugf("sql log : err[%v],time[%.3f],rows[%v]\n%v", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
