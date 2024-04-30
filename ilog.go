package log

import (
	"fmt"

	"github.com/muesli/termenv"
)

type Level int
type Flag int
type writeBuffer struct {
	buffer []byte
}

const (
	LevelFatal Level = 1 << 0
	LevelPanic Level = 1 << 1
	LevelPrint Level = 1 << 2

	LevelInfo  Level = 1 << 3
	LevelWarn  Level = 1 << 4
	LevelError Level = 1 << 5
	LevelDebug Level = 1 << 6
	LevelStack Level = 1 << 7

	LevelAll Level = LevelFatal | LevelPanic | LevelPrint | LevelInfo | LevelWarn | LevelError | LevelDebug | LevelStack
)
const (
	FlagTime      Flag = 1 << 0
	FlagTimeStamp Flag = 1 << 1
	FlagLevel     Flag = 1 << 2
	FlagShortFile Flag = 1 << 3
	FlagLongFile  Flag = 1 << 4
	FlagNewLine   Flag = 1 << 5
	FlagColor     Flag = 1 << 6

	FlagAll Flag = FlagTime | FlagTimeStamp | FlagLevel | FlagShortFile | FlagLongFile | FlagNewLine | FlagColor
)

const (
	DefaultDepth int = 4
)

var (
	lvs = map[Level]string{
		LevelFatal: "F",
		LevelPanic: "P",
		LevelPrint: "R",
		LevelInfo:  "I",
		LevelWarn:  "W",
		LevelError: "E",
		LevelDebug: "D",
		LevelStack: "S",
	}
)

var colors = map[Level][2]func(a ...any) string{
	LevelInfo:  {ANSIColor("", "#8AE234", true, false), ANSIColor("", "#468A05", false, false)},
	LevelWarn:  {ANSIColor("", "#FCE94F", true, false), ANSIColor("", "#C4A000", false, false)},
	LevelError: {ANSIColor("#FFFF55", "#EF2929", true, false), ANSIColor("", "#CC0000", false, false)},
	LevelStack: {ANSIColor("", "#34E2E2", true, false), ANSIColor("", "#046566", false, false)},
	LevelDebug: {ANSIColor("", "#CF269C", true, false), ANSIColor("", "#7B296E", false, false)},
	LevelFatal: {ANSIColor("#FF0000", "#FFFF55", true, true), ANSIColor("", "#FF0000", false, false)},
	LevelPanic: {ANSIColor("#FFFF55", "#FF0000", true, true), ANSIColor("", "#FF0000", false, false)},
	LevelPrint: {ANSIColor("", "#DEDEDC", false, false), ANSIColor("", "#424242", false, false)},
}

type ILogger interface {
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Fatalln(args ...any)

	Panic(args ...any)
	Panicf(format string, args ...any)
	Panicln(args ...any)

	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)

	Info(args ...any)
	Infof(format string, args ...any)
	Infoln(args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)

	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)

	Stack(args ...any)
	Stackf(format string, args ...any)
	Stackln(args ...any)

	Json(lv Level, data any, args ...any)
	LevelRename(lv Level, newName string)
}

type IHandler interface {
	Write(b []byte) (n int, err error)
	Close() error
}

func ANSIColor(background, foreground string, bold, blink bool) func(a ...any) string {
	b := termenv.ANSI.Color(background)
	f := termenv.ANSI.Color(foreground)
	return func(a ...any) string {
		style := termenv.String(fmt.Sprint(a...))
		style = style.Background(b).Foreground(f)
		if bold {
			style = style.Bold()
		}
		if blink {
			style = style.Blink()
		}
		return style.String()
	}
}
