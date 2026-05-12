package log

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	handler     IHandler
	enableColor bool
	shortName   bool
	flagTime    FLAG_TIME
	level       int
	pool        *sync.Pool
}

type writePool struct {
	buffer []byte
}

func poolNew() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return &writePool{
				buffer: make([]byte, 0, 1024),
			}
		},
	}
}

func (l *Logger) withPrefix(lv int, buf *writePool) {
	now := time.Now()
	switch l.flagTime {
	case FLAG_TIME_DATE:
		buf.buffer = append(buf.buffer, now.Format("2006/01/02 ")...)
	case FLAG_TIME_TIME:
		buf.buffer = append(buf.buffer, now.Format("15:04:05.000 ")...)
	case FLAG_TIME_DATETIME:
		buf.buffer = append(buf.buffer, now.Format("2006/01/02 15:04:05.000 ")...)
	case FLAG_TIME_TIMESTAMP:
		buf.buffer = append(buf.buffer, []byte(strconv.Itoa(int(now.UnixMilli())))...)
	case FLAG_TIME_NONE:
		// none time
	}

	attr := LV_ATTRS[lv]
	flagName := ifs(l.shortName, attr.ShortName, attr.Name)
	if l.enableColor {
		buf.buffer = append(buf.buffer, []byte(ColorWrap(flagName, attr.Color...))...)
	} else {
		buf.buffer = append(buf.buffer, []byte(flagName)...)
	}
	buf.buffer = append(buf.buffer, 0x20)

	if lv == LV_DEBUG || lv == LV_ERROR || lv == LV_FATAL || lv == LV_PANIC {
		file, line, fn := WhoCalledMe()
		file = ShortFileName(file)
		if l.enableColor {
			buf.buffer = append(buf.buffer, fmt.Appendf(nil, "%s%s%s%s:%d %s%s%s ", COLOR_CTRL_RESET, COLOR_FG_YELLOW, COLOR_CTRL_UNDERLINE, file, line, COLOR_FG_RED, fn, COLOR_CTRL_RESET)...)
		} else {
			buf.buffer = append(buf.buffer, fmt.Appendf(nil, "%s:%d %s ", file, line, fn)...)
		}
	}
}

func (l *Logger) colorArgs(needSpace bool, args ...any) []any {
	result := make([]any, 0)
	for _, arg := range args {
		result = append(result, l.colorTypes(arg, ""))
		if needSpace {
			result = append(result, " ")
		}
	}
	return result
}

func (l *Logger) colorTypes(arg any, verb string) string {
	str := ifs(verb == "", fmt.Sprint(arg), fmt.Sprintf(verb, arg))
	switch reflect.ValueOf(arg).Kind() {
	case reflect.String:
		return str
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Ptr, reflect.Uintptr,
		reflect.Complex64, reflect.Complex128:
		return ColorWrap(str, COLOR_FG_MAGENTA)
	case reflect.Bool, reflect.Array, reflect.Slice, reflect.Map:
		return ColorWrap(str, COLOR_FG_CYAN)
	case reflect.Struct, reflect.Chan, reflect.Func, reflect.Interface:
		return ColorWrap(str, COLOR_FG_BLUE)
	default:
		return str
	}
}

func (l *Logger) colorFormatArgs(format string, args ...any) string {
	reg := regexp.MustCompile(`%(\[[1-9]\d*\])?([+\-# 0]*)(\d+|\*)?(\.(\d+|\*))?([a-zA-Z%])`)
	matches := reg.FindAllStringSubmatchIndex(format, -1)
	var sb strings.Builder
	var lastIndex, argIndex int = 0, 0
	for _, match := range matches {
		start, end := match[0], match[1]
		verb := format[start:end]
		sb.WriteString(format[lastIndex:start])
		if verb == "%%" {
			sb.WriteString("%")
			lastIndex = end
			continue
		}

		var idx int
		hasIndex := false
		m2, m3 := match[2], match[3]
		if m2 > 0 && m3 > 0 {
			idxStr := format[m2+1 : m3-1]
			idx, _ = strconv.Atoi(idxStr)
			hasIndex = true
		}
		var arg any
		if hasIndex {
			if idx < 1 || idx > len(args) {
				arg = nil
			} else {
				arg = args[idx-1]
			}
		} else {
			if argIndex >= len(args) {
				arg = nil
			} else {
				arg = args[argIndex]
				argIndex++
			}
		}
		sb.WriteString(l.colorTypes(arg, verb))
		lastIndex = end
	}
	sb.WriteString(format[lastIndex:])
	return sb.String()
}

func (l *Logger) log(lv int, args ...any) {
	if lv < l.level {
		return
	}
	buf := l.pool.Get().(*writePool)
	buf.buffer = buf.buffer[:0]
	defer l.pool.Put(buf)
	l.withPrefix(lv, buf)
	if l.enableColor {
		args = l.colorArgs(true, args...)
	}
	buf.buffer = append(buf.buffer, fmt.Sprint(args...)...)
	l.handler.Write(buf.buffer)
	l.handler.Write([]byte("\n"))
}

func (l *Logger) logf(lv int, format string, args ...any) {
	if lv < l.level {
		return
	}
	buf := l.pool.Get().(*writePool)
	buf.buffer = buf.buffer[:0]
	defer l.pool.Put(buf)
	l.withPrefix(lv, buf)
	if l.enableColor {
		buf.buffer = append(buf.buffer, l.colorFormatArgs(format, args...)...)
	} else {
		buf.buffer = append(buf.buffer, fmt.Sprintf(format, args...)...)
	}
	l.handler.Write(buf.buffer)
	l.handler.Write([]byte("\n"))
}

func (l *Logger) Fatal(args ...any) {
	l.log(LV_FATAL, args...)
	os.Exit(0)
}
func (l *Logger) Fatalf(format string, args ...any) {
	l.logf(LV_FATAL, format, args...)
	os.Exit(0)
}
func (l *Logger) Fatalln(args ...any) {
	l.log(LV_FATAL, args...)
	os.Exit(0)
}

func (l *Logger) Panic(args ...any) {
	l.log(LV_PANIC, args...)
	panic(struct{}{})
}
func (l *Logger) Panicf(format string, args ...any) {
	l.logf(LV_PANIC, format, args...)
	panic(struct{}{})
}
func (l *Logger) Panicln(args ...any) {
	l.log(LV_PANIC, args...)
	panic(struct{}{})
}

func (l *Logger) Print(args ...any) {
	l.log(LV_PRINT, args...)
}
func (l *Logger) Printf(format string, args ...any) {
	l.logf(LV_PRINT, format, args...)
}
func (l *Logger) Println(args ...any) {
	l.log(LV_PRINT, args...)
}

func (l *Logger) Info(args ...any) {
	l.log(LV_INFO, args...)
}
func (l *Logger) Infof(format string, args ...any) {
	l.logf(LV_INFO, format, args...)
}
func (l *Logger) Infoln(args ...any) {
	l.log(LV_INFO, args...)
}

func (l *Logger) Warn(args ...any) {
	l.log(LV_WARN, args...)
}
func (l *Logger) Warnf(format string, args ...any) {
	l.logf(LV_WARN, format, args...)
}
func (l *Logger) Warnln(args ...any) {
	l.log(LV_WARN, args...)
}

func (l *Logger) Error(args ...any) {
	l.log(LV_ERROR, args...)
}
func (l *Logger) Errorf(format string, args ...any) {
	l.logf(LV_ERROR, format, args...)
}
func (l *Logger) Errorln(args ...any) {
	l.log(LV_ERROR, args...)
}

func (l *Logger) Debug(args ...any) {
	l.log(LV_DEBUG, args...)
}
func (l *Logger) Debugf(format string, args ...any) {
	l.logf(LV_DEBUG, format, args...)
}
func (l *Logger) Debugln(args ...any) {
	l.log(LV_DEBUG, args...)
}
