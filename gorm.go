package log

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	glogger "gorm.io/gorm/logger"
)

type GormLoggerHandler struct {
	handler                   IHandler
	LogLevel                  glogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	ColorFul                  bool
}

func NewGormLoggerHandler(handler IHandler) *GormLoggerHandler {
	return &GormLoggerHandler{
		handler: handler,
	}
}

func (that *GormLoggerHandler) Write(b []byte) (n int, err error) {
	return that.handler.Write(b)
}
func (that *GormLoggerHandler) Close() error {
	return that.handler.Close()
}

func (that *GormLoggerHandler) LogMode(level glogger.LogLevel) glogger.Interface {
	that.LogLevel = level
	return that
}
func (that *GormLoggerHandler) SetSlowThreshold(value time.Duration) glogger.Interface {
	that.SlowThreshold = value
	return that
}
func (that *GormLoggerHandler) SetIgnoreRecordNotFoundError(value bool) glogger.Interface {
	that.IgnoreRecordNotFoundError = value
	return that
}
func (that *GormLoggerHandler) SetColorful(value bool) glogger.Interface {
	that.ColorFul = value
	return that
}

func (that *GormLoggerHandler) format(level Level, format string, args ...any) (date, lv, msg string) {
	date = time.Now().Format("2006/01/02 15:04:05.000")
	lv = fmt.Sprintf("[%s]", lvs[level])
	msg = fmt.Sprintf(format, args...)
	if that.ColorFul {
		date = colors[LevelInfo][0](date)
		lv = colors[LevelInfo][0](lv)
		msg = colors[LevelInfo][1](msg)
	}
	return
}

func (that *GormLoggerHandler) Info(ctx context.Context, format string, args ...any) {
	if that.LogLevel >= glogger.Info {
		date, lv, msg := that.format(LevelInfo, format, args...)
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s\n", date, lv, msg)))
	}
}

func (that *GormLoggerHandler) Warn(ctx context.Context, format string, args ...any) {
	if that.LogLevel >= glogger.Warn {
		date, lv, msg := that.format(LevelWarn, format, args...)
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s\n", date, lv, msg)))
	}
}
func (that *GormLoggerHandler) Error(ctx context.Context, format string, args ...any) {
	if that.LogLevel >= glogger.Error {
		date, lv, msg := that.format(LevelError, format, args...)
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s\n", date, lv, msg)))
	}
}

func (that *GormLoggerHandler) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if that.LogLevel <= glogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	file, line := _caller_file_line()
	file__line := fmt.Sprintf("%s:%d", file, line)
	rows_times := fmt.Sprintf("[rows:%d %.3fms]", rows, float64(elapsed.Nanoseconds())/1e6)
	date, lv, _ := that.format(LevelDebug, "")
	slow := "[SLOW]"
	if that.ColorFul {
		file__line = colors[LevelStack][1](file__line)
		rows_times = colors[LevelInfo][1](rows_times)
		sql = colors[LevelPrint][1](sql)
		slow = colors[LevelError][1](slow)
	}
	switch {
	case err != nil && that.LogLevel >= glogger.Error && (!errors.Is(err, glogger.ErrRecordNotFound) || !that.IgnoreRecordNotFoundError):
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s %s %s \nError: %v\n", date, lv, file__line, rows_times, sql, err)))
	case elapsed > that.SlowThreshold && that.SlowThreshold != 0 && that.LogLevel >= glogger.Warn:
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s %s %s %s\n", date, lv, file__line, rows_times, slow, sql)))
	case that.LogLevel == glogger.Info:
		that.handler.Write([]byte(fmt.Sprintf("%s %s %s %s %s\n", date, lv, file__line, rows_times, sql)))
	}
}

func _caller_file_line() (string, int) {
	pcs := [13]uintptr{}
	len := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:len])
	for i := 0; i < len; i++ {
		frame, _ := frames.Next()
		if !strings.Contains(frame.File, "/gorm.io/") {
			return frame.File, frame.Line
		}
	}
	return "", 0
}
