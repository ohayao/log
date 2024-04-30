package log

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	_default *Logger
)

func init() {
	handler := NewStreamHandler(os.Stderr)
	_default = NewLogger(handler)
	SetFlags(FlagAll &^ FlagTimeStamp &^ FlagNewLine)
	SetLevels(LevelAll)
}

func Fatal(args ...any) {
	_default.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	_default.Fatalf(format, args...)
}

func Fatalln(args ...any) {
	_default.Fatalln(args...)
}

func Panic(args ...any) {
	_default.Panic(args...)
}

func Panicf(format string, args ...any) {
	_default.Panicf(format, args...)
}

func Panicln(args ...any) {
	_default.Panicln(args...)
}

func Print(args ...any) {
	_default.Print(args...)
}

func Printf(format string, args ...any) {
	_default.Printf(format, args...)
}

func Println(args ...any) {
	_default.Println(args...)
}

func Info(args ...any) {
	_default.Info(args...)
}

func Infof(format string, args ...any) {
	_default.Infof(format, args...)
}

func Infoln(args ...any) {
	_default.Infoln(args...)
}

func Warn(args ...any) {
	_default.Warn(args...)
}

func Warnf(format string, args ...any) {
	_default.Warnf(format, args...)
}

func Warnln(args ...any) {
	_default.Warnln(args...)
}

func Error(args ...any) {
	_default.Error(args...)
}

func Errorf(format string, args ...any) {
	_default.Errorf(format, args...)
}

func Errorln(args ...any) {
	_default.Errorln(args...)
}

func Debug(args ...any) {
	_default.Debug(args...)
}

func Debugf(format string, args ...any) {
	_default.Debugf(format, args...)
}

func Debugln(args ...any) {
	_default.Debugln(args...)
}

func Stack(args ...any) {
	_default._stackln(DefaultDepth+1, fmt.Sprint(args...))
}

func Stackf(format string, args ...any) {
	_default._stackln(DefaultDepth+1, fmt.Sprintf(format, args...))
}

func Stackln(args ...any) {
	_default._stackln(DefaultDepth+1, args...)
}

func Json(lv Level, data any, args ...any) {
	if lv == LevelStack {
		bs, _ := json.Marshal(data)
		str := fmt.Sprint(fmt.Sprint(args...), string(bs))
		_default._stackln(DefaultDepth+1, str)
	} else {
		_default.Json(lv, data, args...)

	}
}

func LevelRename(lv Level, newName string) {
	_default.LevelRename(lv, newName)
}

func SetLevels(levels ...Level) {
	_default.SetLevels(levels...)
}

func SetFlags(flgas ...Flag) {
	_default.SetFlags(flgas...)
}
