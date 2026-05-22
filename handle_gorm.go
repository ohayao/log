package log

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	base                      *Logger
	level                     logger.LogLevel
	_level                    int
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.level = level
	switch level {
	case logger.Silent:
		g._level = LV_DEBUG - 1
	case logger.Error:
		g._level = LV_ERROR
	case logger.Warn:
		g._level = LV_WARN
	case logger.Info:
		g._level = LV_INFO
	}
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if LV_INFO >= g._level {
		g.base.logf_gorm(LV_INFO, msg, data...)
	}
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if LV_WARN >= g._level {
		g.base.logf_gorm(LV_WARN, msg, data...)
	}
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if LV_ERROR >= g._level {
		g.base.logf_gorm(LV_ERROR, msg, data...)
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.level == logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && LV_ERROR >= g._level && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.ignoreRecordNotFoundError):
		sql, rows := fc()
		if g.base.enableColor {
			sql = ColorWrap(sql, COLOR_FG_CYAN)
		}
		file, line := _caller_file_line()
		file = ShortFileName(file)
		g.base.logf_gorm(LV_ERROR, "[%s:%d rows:%d %.3fms] %s err: %v", file, line, rows, float64(elapsed.Nanoseconds())/1e6, sql, err)
	case elapsed >= g.slowThreshold && LV_WARN >= g._level && g.slowThreshold > 0:
		sql, rows := fc()
		if g.base.enableColor {
			sql = ColorWrap(sql, COLOR_FG_CYAN)
		}
		file, line := _caller_file_line()
		file = ShortFileName(file)
		g.base.logf_gorm(LV_WARN, "[%s:%d rows:%d %.3fms] %s", file, line, rows, float64(elapsed.Nanoseconds())/1e6, sql)
	case LV_INFO >= g._level:
		sql, rows := fc()
		if g.base.enableColor {
			sql = ColorWrap(sql, COLOR_FG_CYAN)
		}
		file, line := _caller_file_line()
		file = ShortFileName(file)
		g.base.logf_gorm(LV_INFO, "[%s:%d rows:%d %.3fms] %s", file, line, rows, float64(elapsed.Nanoseconds())/1e6, sql)
	}
}

func _caller_file_line() (string, int) {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if strings.HasPrefix(frame.Function, packagePrefix) ||
			strings.HasPrefix(frame.Function, "gorm.io/") {
			if !more {
				break
			}
			continue
		}
		return frame.File, frame.Line
	}
	return "", 0
}

func NewGormLogger(baseLogger *Logger, ignoreRecordNotFoundError bool, slowThreshold time.Duration) *GormLogger {
	handler := &GormLogger{
		base:                      baseLogger,
		ignoreRecordNotFoundError: ignoreRecordNotFoundError,
		slowThreshold:             slowThreshold,
		level:                     logger.Info,
		_level:                    LV_INFO,
	}
	return handler
}
