package log

type (
	FLAG_TIME int

	COLOR_ENUM string

	lvAttr struct {
		Name, ShortName string
		Color           []COLOR_ENUM
	}
)

const (
	// 不显示时间
	FLAG_TIME_NONE FLAG_TIME = 0
	// 日期: eg: 2022/01/01
	FLAG_TIME_DATE FLAG_TIME = 1
	// 时间: eg: 15:04:05.000
	FLAG_TIME_TIME FLAG_TIME = 2
	// 日期时间: eg: 2022/01/01 15:04:05.000
	FLAG_TIME_DATETIME FLAG_TIME = 3
	// 时间戳
	FLAG_TIME_TIMESTAMP FLAG_TIME = 4
)

const (
	LV_DEBUG = iota
	LV_PRINT
	LV_INFO
	LV_WARN
	LV_ERROR
	LV_PANIC
	LV_FATAL
)

const (
	COLOR_CTRL_RESET     COLOR_ENUM = "\x1b[0m"
	COLOR_CTRL_BOLD      COLOR_ENUM = "\x1b[1m"
	COLOR_CTRL_UNDERLINE COLOR_ENUM = "\x1b[4m"
	COLOR_CTRL_FLASH     COLOR_ENUM = "\x1b[5m"
	COLOR_CTRL_REVERSE   COLOR_ENUM = "\x1b[7m"

	COLOR_FG_BLACK   COLOR_ENUM = "\x1b[30m"
	COLOR_FG_RED     COLOR_ENUM = "\x1b[31m"
	COLOR_FG_GREEN   COLOR_ENUM = "\x1b[32m"
	COLOR_FG_YELLOW  COLOR_ENUM = "\x1b[33m"
	COLOR_FG_BLUE    COLOR_ENUM = "\x1b[34m"
	COLOR_FG_MAGENTA COLOR_ENUM = "\x1b[35m"
	COLOR_FG_CYAN    COLOR_ENUM = "\x1b[36m"
	COLOR_FG_WHITE   COLOR_ENUM = "\x1b[37m"

	COLOR_BG_BLACK   COLOR_ENUM = "\x1b[40m"
	COLOR_BG_RED     COLOR_ENUM = "\x1b[41m"
	COLOR_BG_GREEN   COLOR_ENUM = "\x1b[42m"
	COLOR_BG_YELLOW  COLOR_ENUM = "\x1b[43m"
	COLOR_BG_BLUE    COLOR_ENUM = "\x1b[44m"
	COLOR_BG_MAGENTA COLOR_ENUM = "\x1b[45m"
	COLOR_BG_CYAN    COLOR_ENUM = "\x1b[46m"
	COLOR_BG_WHITE   COLOR_ENUM = "\x1b[47m"
)

var (
	LV_ATTRS = map[int]lvAttr{
		LV_DEBUG: {Name: "DBG", ShortName: "D", Color: []COLOR_ENUM{COLOR_CTRL_RESET, COLOR_CTRL_BOLD}},
		LV_PRINT: {Name: "PRT", ShortName: "P", Color: []COLOR_ENUM{COLOR_FG_CYAN, COLOR_CTRL_BOLD}},
		LV_INFO:  {Name: "INF", ShortName: "I", Color: []COLOR_ENUM{COLOR_FG_BLUE, COLOR_CTRL_BOLD}},
		LV_WARN:  {Name: "WRN", ShortName: "W", Color: []COLOR_ENUM{COLOR_FG_YELLOW, COLOR_CTRL_BOLD}},
		LV_ERROR: {Name: "ERR", ShortName: "E", Color: []COLOR_ENUM{COLOR_FG_RED, COLOR_CTRL_BOLD}},
		LV_PANIC: {Name: "PNC", ShortName: "S", Color: []COLOR_ENUM{COLOR_FG_MAGENTA, COLOR_CTRL_BOLD}},
		LV_FATAL: {Name: "FAT", ShortName: "F", Color: []COLOR_ENUM{COLOR_FG_MAGENTA, COLOR_CTRL_BOLD}},
	}
)

type (
	ILogger interface {
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
	}

	IHandler interface {
		Write(b []byte) (n int, err error)
		Close() (err error)
	}
)

func ifs[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
