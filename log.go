package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
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
// eg: SetFlags(FlagTime,FlagLevel,FlagColor).
// eg: SetFlags(FlagAll &^ FlagNewLine) 所有标记，除去换行标记.
// 当Level标记中未排除Stack,Fatal,Panic时，且未标记长文件或短文件时,
// 自动添加一个短文件标记
func (that *Logger) SetFlags(flags ...Flag) {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.flag = 0
	for _, f := range flags {
		that.flag = that.flag | int(f)
	}
	if !(that.flag&int(FlagShortFile) > 0 || that.flag&int(FlagLongFile) > 0) && (that.level&int(LevelStack) > 0 || that.level&int(LevelFatal) > 0 || that.level&int(LevelPanic) > 0) {
		that.flag |= int(FlagShortFile)
	}
}

// SetLevels
// eg: SetLevels(LevelInfo,LevelError,LevelStack).
// eg: SetLevels(Level &^ LevelDebug) 所有标记，除去调试信息.
// 当Level中包含Stack,Fatal,Panic时，且flag未有文件标记时,
// 自动在flag中添加短文件标记
func (that *Logger) SetLevels(levels ...Level) {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.level = 0
	for _, lv := range levels {
		that.level = that.level | int(lv)
	}
	if (that.level&int(LevelStack) > 0 || that.level&int(LevelFatal) > 0 || that.level&int(LevelPanic) > 0) && !(that.flag&int(FlagShortFile) > 0 || that.flag&int(FlagLongFile) > 0) {
		that.flag |= int(FlagShortFile)
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
	if that.level&int(lv) < 1 {
		return
	}
	buf := that.buff.Get().(*writeBuffer)
	buf.buffer = buf.buffer[0:0]
	defer that.buff.Put(buf)

	buf.buffer = append(buf.buffer, that.prefix(DefaultDepth, lv).Bytes()...)
	if that.flag&int(FlagColor) > 0 {
		buf.buffer = append(buf.buffer, colors[lv][1](data)...)
	} else {
		buf.buffer = append(buf.buffer, data...)
	}

	if len(buf.buffer) > 0 && buf.buffer[len(buf.buffer)-1] != '\n' {
		buf.buffer = append(buf.buffer, '\n')
	}

	that.lock.Lock()
	that.handler.Write(buf.buffer)
	that.lock.Unlock()
}

func (that *Logger) writeStack(depthStart int, lv Level, data string) {
	if that.level&int(lv) < 1 {
		return
	}
	buf := that.buff.Get().(*writeBuffer)
	buf.buffer = buf.buffer[0:0]
	defer that.buff.Put(buf)

	if depthStart < DefaultDepth {
		depthStart = DefaultDepth
	}

	buf.buffer = append(buf.buffer, that.prefix(depthStart, lv).Bytes()...)
	if that.flag&int(FlagColor) > 0 {
		buf.buffer = append(buf.buffer, colors[lv][1](data)...)
	} else {
		buf.buffer = append(buf.buffer, data...)
	}

	if len(buf.buffer) > 0 && buf.buffer[len(buf.buffer)-1] != '\n' {
		buf.buffer = append(buf.buffer, '\n')
	}

	that.lock.Lock()
	that.handler.Write(buf.buffer)
	that.lock.Unlock()
}

func (that *Logger) prefix(depthStart int, lv Level) *bytes.Buffer {
	buf := bytes.NewBufferString("")
	now := time.Now()
	if that.flag&int(FlagTime) > 0 {
		buf.WriteString(now.Format("2006/01/02 15:04:05.000 "))
	}

	if that.flag&int(FlagTimeStamp) > 0 {
		buf.WriteString(fmt.Sprintf("%d ", now.UnixMicro()/1e3))
	}

	if that.flag&int(FlagLevel) > 0 {
		buf.WriteString(fmt.Sprintf("[%s] ", lvs[lv]))
	}

	if a, b := that.flag&int(FlagShortFile), that.flag&int(FlagLongFile); (lv == LevelStack || lv == LevelFatal || lv == LevelPanic) && (a > 0 || b > 0) {
		buf.WriteByte('[')
		that._stack(depthStart, a > 0, buf)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}

	if that.flag&int(FlagColor) > 0 && buf.Len() > 1 {
		str := colors[lv][0](string(buf.Bytes()[:buf.Len()-1]))
		buf.Reset()
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	if that.flag&int(FlagNewLine) > 0 {
		buf.WriteByte('\n')
	}
	return buf
}

func (that *Logger) _stack(skip int, isShort bool, buf *bytes.Buffer) {
	list := make([]string, 0)
	for {
		if pc, file, line, ok := runtime.Caller(skip); ok {
			fn := runtime.FuncForPC(pc)
			if fn == nil || strings.HasPrefix(strings.ToLower(fn.Name()), "runtime.") {
				break
			} else {
				if isShort {
					list = append(list, fmt.Sprintf("%s:%d", _shortFile(file), line))
				} else {
					list = append(list, fmt.Sprintf("%s:%d", _longFile(file), line))
				}
			}
		} else {
			break
		}
		skip++
	}

	buf.WriteString(strings.Join(list, " <- "))
}

func _shortFile(file string) string {
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			return file[i+1:]
		}
	}
	return file
}

func _longFile(file string) string {
	file = strings.ReplaceAll(file, "\\", "/")
	if dir, err := os.Getwd(); err != nil {
		return file
	} else {
		dir := strings.ReplaceAll(dir, "\\", "/")
		return strings.TrimPrefix(strings.TrimPrefix(file, dir), "/")
	}
}

func (that *Logger) Close() {
	that.lock.Lock()
	defer that.lock.Unlock()
	that.handler.Close()
}

func (that *Logger) Fatal(args ...any) {
	that.writeStack(DefaultDepth, LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}
func (that *Logger) Fatalf(format string, args ...any) {
	that.writeStack(DefaultDepth, LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}
func (that *Logger) Fatalln(args ...any) {
	that.writeStack(DefaultDepth, LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func (that *Logger) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	that.writeStack(DefaultDepth, LevelPanic, fmt.Sprint(args...))
	panic(msg)
}
func (that *Logger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	that.writeStack(DefaultDepth, LevelPanic, fmt.Sprint(args...))
	panic(msg)
}
func (that *Logger) Panicln(args ...any) {
	msg := fmt.Sprint(args...)
	that.writeStack(DefaultDepth, LevelPanic, fmt.Sprint(args...))
	panic(msg)
}

func (that *Logger) Print(args ...any) {
	that.write(LevelPrint, fmt.Sprint(args...))
}
func (that *Logger) Printf(format string, args ...any) {
	that.write(LevelPrint, fmt.Sprintf(format, args...))
}
func (that *Logger) Println(args ...any) {
	that.write(LevelPrint, fmt.Sprint(args...))
}

func (that *Logger) Info(args ...any) {
	that.write(LevelInfo, fmt.Sprint(args...))
}
func (that *Logger) Infof(format string, args ...any) {
	that.write(LevelInfo, fmt.Sprintf(format, args...))
}
func (that *Logger) Infoln(args ...any) {
	that.write(LevelInfo, fmt.Sprint(args...))
}

func (that *Logger) Warn(args ...any) {
	that.write(LevelWarn, fmt.Sprint(args...))
}
func (that *Logger) Warnf(format string, args ...any) {
	that.write(LevelWarn, fmt.Sprintf(format, args...))
}
func (that *Logger) Warnln(args ...any) {
	that.write(LevelWarn, fmt.Sprint(args...))
}

func (that *Logger) Error(args ...any) {
	that.write(LevelError, fmt.Sprint(args...))
}
func (that *Logger) Errorf(format string, args ...any) {
	that.write(LevelError, fmt.Sprintf(format, args...))
}
func (that *Logger) Errorln(args ...any) {
	that.write(LevelError, fmt.Sprint(args...))
}

func (that *Logger) Debug(args ...any) {
	that.write(LevelDebug, fmt.Sprint(args...))
}
func (that *Logger) Debugf(format string, args ...any) {
	that.write(LevelDebug, fmt.Sprintf(format, args...))
}
func (that *Logger) Debugln(args ...any) {
	that.write(LevelDebug, fmt.Sprint(args...))
}

// Stack
func (that *Logger) Stack(args ...any) {
	that.writeStack(DefaultDepth, LevelStack, fmt.Sprint(args...))
}

// Stackf
func (that *Logger) Stackf(format string, args ...any) {
	that.writeStack(DefaultDepth, LevelStack, fmt.Sprintf(format, args...))
}

// Stackln
func (that *Logger) Stackln(args ...any) {
	that.writeStack(DefaultDepth, LevelStack, fmt.Sprint(args...))
}

func (that *Logger) _stackln(depth int, args ...any) {
	that.writeStack(depth, LevelStack, fmt.Sprint(args...))
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
