package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	handler IHandler
	flag    int
	level   int
	buff    sync.Pool
	lock    sync.Mutex
}

func NewLogger(handler IHandler) *Logger {
	var log = new(Logger)
	log.handler = handler
	log.buff = sync.Pool{
		New: func() interface{} {
			buffer := &writeBuffer{buffer: make([]byte, 0, 1024)}
			return buffer
		},
	}
	return log
}

// SetFlags
func (that *Logger) SetFlags(flags ...Flag) {
	that.lock.Lock()
	defer that.lock.Unlock()
	for _, f := range flags {
		that.flag = that.flag | int(f)
	}
}

// SetLevels
func (that *Logger) SetLevels(levels ...Level) {
	that.lock.Lock()
	defer that.lock.Unlock()
	for _, lv := range levels {
		that.level = that.level | int(lv)
	}
}

// LevelRename
func (that *Logger) LevelRename(lv Level, newName string) {
	that.lock.Lock()
	defer that.lock.Unlock()
	lvs[lv] = newName
}

// write
func (that *Logger) write(lv Level, data string) {
	buf := that.buff.Get().(*writeBuffer)
	buf.buffer = buf.buffer[0:0]
	defer that.buff.Put(buf)

	buf.buffer = append(buf.buffer, that.prefix(DefaultDepth, DefaultDepth, lv).Bytes()...)
	buf.buffer = append(buf.buffer, data...)

	if buf.buffer[len(buf.buffer)-1] != '\n' {
		buf.buffer = append(buf.buffer, '\n')
	}

	that.lock.Lock()
	that.handler.Write(buf.buffer)
	that.lock.Unlock()
}

func (that *Logger) writeStack(depthStart, depthEnd int, lv Level, data string) {
	buf := that.buff.Get().(*writeBuffer)
	buf.buffer = buf.buffer[0:0]
	defer that.buff.Put(buf)

	if depthStart < DefaultDepth {
		depthStart = DefaultDepth
	}
	if depthEnd < depthStart {
		depthEnd = depthStart
	}

	buf.buffer = append(buf.buffer, that.prefix(depthStart, depthEnd, lv).Bytes()...)
	buf.buffer = append(buf.buffer, data...)

	if buf.buffer[len(buf.buffer)-1] != '\n' {
		buf.buffer = append(buf.buffer, '\n')
	}

	that.lock.Lock()
	that.handler.Write(buf.buffer)
	that.lock.Unlock()
}

func (that *Logger) prefix(depthStart, depthEnd int, lv Level) *bytes.Buffer {
	buf := bytes.NewBufferString("")

	if that.flag&int(FlagTime) > 0 {
		buf.WriteString(time.Now().Format("2006/01/02 15:04:05.000 "))
	}

	if that.flag&int(FlagTimeStamp) > 0 {
		buf.WriteString(fmt.Sprintf("%d ", time.Now().UnixMicro()/1e3))
	}

	if that.flag&int(FlagLevel) > 0 {
		lv_str := fmt.Sprintf("[%s]", lvs[lv])
		if that.flag&int(FlagColor) > 0 {
			lv_str = colors[lv](lv_str)
		}
		buf.WriteString(lv_str)
		buf.WriteByte(' ')
	}

	if a, b := that.flag&int(FlagShortFile), that.flag&int(FlagLongFile); lv == LevelStack && (a > 0 || b > 0) {
		buf.WriteByte('[')
		that._stack(depthStart, depthEnd, a > 0, buf)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}

	if that.flag&int(FlagNewLine) > 0 {
		buf.WriteByte('\n')
	}
	return buf
}

func (that *Logger) _stack(start, end int, isShort bool, buf *bytes.Buffer) {
	var last_file, last_line = "???", 0
	for x := start; x <= end; x++ {
		if _, file, line, ok := runtime.Caller(x); ok {
			if last_file == file && last_line == line {
				break
			} else {
				if x > start {
					buf.WriteString(" > ")
				}
				if isShort {
					buf.WriteString(fmt.Sprintf("%s:%d", _shortFile(file), line))
				} else {
					buf.WriteString(fmt.Sprintf("%s:%d", file, line))
				}
			}
			last_file = file
			last_line = line
		} else {
			break
		}
	}
}

func _shortFile(file string) string {
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			return file[i+1:]
		}
	}
	return file
}

func (that *Logger) Close() {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.handler.Close()
}

func (that *Logger) Fatal(args ...any) {
	that.write(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}
func (that *Logger) Fatalf(format string, args ...any) {
	that.write(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}
func (that *Logger) Fatalln(args ...any) {
	that.write(LevelFatal, fmt.Sprintln(args...))
	os.Exit(1)
}

func (that *Logger) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	that.write(LevelPanic, msg)
	panic(msg)
}
func (that *Logger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	that.write(LevelPanic, msg)
	panic(msg)
}
func (that *Logger) Panicln(args ...any) {
	msg := fmt.Sprintln(args...)
	that.write(LevelPanic, msg)
	panic(msg)
}

func (that *Logger) Print(args ...any) {
	that.write(LevelPrint, fmt.Sprint(args...))
}
func (that *Logger) Printf(format string, args ...any) {
	that.write(LevelPrint, fmt.Sprintf(format, args...))
}
func (that *Logger) Println(args ...any) {
	that.write(LevelPrint, fmt.Sprintln(args...))
}

func (that *Logger) Info(args ...any) {
	that.write(LevelInfo, fmt.Sprint(args...))
}
func (that *Logger) Infof(format string, args ...any) {
	that.write(LevelInfo, fmt.Sprintf(format, args...))
}
func (that *Logger) Infoln(args ...any) {
	that.write(LevelInfo, fmt.Sprintln(args...))
}

func (that *Logger) Warn(args ...any) {
	that.write(LevelWarn, fmt.Sprint(args...))
}
func (that *Logger) Warnf(format string, args ...any) {
	that.write(LevelWarn, fmt.Sprintf(format, args...))
}
func (that *Logger) Warnln(args ...any) {
	that.write(LevelWarn, fmt.Sprintln(args...))
}

func (that *Logger) Error(args ...any) {
	that.write(LevelError, fmt.Sprint(args...))
}
func (that *Logger) Errorf(format string, args ...any) {
	that.write(LevelError, fmt.Sprintf(format, args...))
}
func (that *Logger) Errorln(args ...any) {
	that.write(LevelError, fmt.Sprintln(args...))
}

func (that *Logger) Debug(args ...any) {
	that.write(LevelDebug, fmt.Sprint(args...))
}
func (that *Logger) Debugf(format string, args ...any) {
	that.write(LevelDebug, fmt.Sprintf(format, args...))
}
func (that *Logger) Debugln(args ...any) {
	that.write(LevelDebug, fmt.Sprintln(args...))
}

// Stack
// depth start from 4
func (that *Logger) Stack(depth int, args ...any) {
	that.writeStack(DefaultDepth, depth, LevelStack, fmt.Sprint(args...))
}

// Stackf
// depth start from 4
func (that *Logger) Stackf(depth int, format string, args ...any) {
	that.writeStack(DefaultDepth, depth, LevelStack, fmt.Sprintf(format, args...))
}

// Stackln
// depth start from 4
func (that *Logger) Stackln(depth int, args ...any) {
	that.writeStack(DefaultDepth, depth, LevelStack, fmt.Sprintln(args...))
}

func (that *Logger) _stackln(depth int, args ...any) {
	that.writeStack(DefaultDepth+1, depth, LevelStack, fmt.Sprintln(args...))
}

func (that *Logger) Json(lv Level, data any, args ...any) {
	bs, _ := json.Marshal(data)
	str := fmt.Sprint(fmt.Sprint(args...), string(bs))
	switch lv {
	case LevelFatal:
		that.Fatalln(str)
	case LevelPanic:
		that.Panicln(str)
	case LevelPrint:
		that.Println(str)
	case LevelInfo:
		that.Infoln(str)
	case LevelWarn:
		that.Warnln(str)
	case LevelError:
		that.Errorln(str)
	case LevelDebug:
		that.Debugln(str)
	case LevelStack:
		that._stackln(DefaultDepth+1, str)
	default:
		that.Infoln(str)
	}
}
