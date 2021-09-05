package dblogger

import (
	"context"
	"fmt"
	"strings"
	"time"

	glogger "gorm.io/gorm/logger"
)

const (
	SQL_PRINT_MAX_LEN = 4096 * 10
)

type DbLogger struct {
	LogLevel glogger.LogLevel
}

func (l *DbLogger) LogMode(glogger.LogLevel) glogger.Interface {
	return l
}

func (l DbLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= glogger.Info {
		Log.Infof(format, args...)
	}
}
func (l DbLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= glogger.Warn {
		Log.Warnf(format, args...)
	}
}
func (l DbLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= glogger.Error {
		Log.Errorf(format, args...)
	}
}
func (l DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == glogger.Info {
		elapsed := time.Since(begin)
		sql, rows := fc()

		prow := fmt.Sprintf("%v", rows)
		if rows == -1 {
			prow = "-"
		}

		// 过滤部分sql
		if l.sqlFilter(sql) {
			return
		}

		// sql 过长会截断 40k 长度一般sql足够
		if len(sql) > SQL_PRINT_MAX_LEN {
			sql = sql[:SQL_PRINT_MAX_LEN] + "...[truncated]"
		}

		if err == nil {
			Log.Debugf("<SQL> : [%v]\ntime : [%.3f]\nrows : [%v]\n", sql, float64(elapsed.Nanoseconds())/1e6, prow)
		} else {
			Log.Debugf("<SQL> : [%v]\ntime : [%.3f]\nrows : [%v]\nerror : [%v],", sql, float64(elapsed.Nanoseconds())/1e6, prow, err)
		}
	}
}

func (l DbLogger) sqlFilter(sql string) bool {
	// no print 标记不打印
	if strings.Contains(sql, "/*no print*/") {
		return true
	}

	if strings.Contains(sql, "information_schema") {
		return true
	}

	if strings.Contains(sql, "SELECT DATABASE()") {
		return true
	}

	if strings.Contains(sql, "sqlite_master") {
		return true
	}

	return false
}

// 可以用于临时不打印日志
type NoLogger struct {
	LogLevel glogger.LogLevel
}

func (l *NoLogger) LogMode(glogger.LogLevel) glogger.Interface {
	return l
}

func (l NoLogger) Info(ctx context.Context, format string, args ...interface{}) {
}
func (l NoLogger) Warn(ctx context.Context, format string, args ...interface{}) {
}
func (l NoLogger) Error(ctx context.Context, format string, args ...interface{}) {
}
func (l NoLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
}
