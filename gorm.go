package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	glogger "gorm.io/gorm/logger"
)

type GormLogHandler struct {
	ctx    context.Context
	Logger *Logger
}

func WrapLoggerForGorm(ctx context.Context, logger *Logger) *GormLogHandler {
	return &GormLogHandler{
		ctx:    ctx,
		Logger: logger,
	}
}

func (that *GormLogHandler) LogMode(level glogger.LogLevel) glogger.Interface {
	switch level {
	case glogger.Info:
		that.Logger.SetLevels(LevelInfo, LevelWarn, LevelError, LevelDebug)
	case glogger.Warn:
		that.Logger.SetLevels(LevelWarn, LevelError, LevelDebug)
	case glogger.Error:
		that.Logger.SetLevels(LevelError, LevelDebug)
	case glogger.Silent:
		that.Logger.SetLevels(LevelDebug) // 禁用日志, 但是debug级别仍有效
	}
	return that
}

func (that *GormLogHandler) Info(ctx context.Context, format string, args ...any) {
	that.Logger.Infof(format, args...)
}
func (that *GormLogHandler) Warn(ctx context.Context, format string, args ...any) {
	that.Logger.Warnf(format, args...)
}
func (that *GormLogHandler) Error(ctx context.Context, format string, args ...any) {
	that.Logger.Errorf(format, args...)
}
func (that *GormLogHandler) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	file, line := _FileWithLineNum()
	if that.Logger.flag&int(FlagShortFile) > 0 {
		file = _shortFile(file)
	}
	file__line := fmt.Sprintf("%s:%d", file, line)
	rows_times := fmt.Sprintf("[rows:%d %.3fms]", rows, float64(elapsed.Nanoseconds())/1e6)
	if that.Logger.flag&int(FlagColor) > 0 {
		file__line = colors[LevelStack][1](file__line)
		rows_times = colors[LevelInfo][1](rows_times)
		sql = colors[LevelPrint][1](sql)
	}
	that.Logger.Debugf("%s %s %s", file__line, rows_times, sql)
}

func _FileWithLineNum() (string, int) {
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
