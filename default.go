package log

var DEFAULT *Logger

func init() {
	DEFAULT = New(newTerminalHandler(nil), WithColor(true), WithMinLevel(LV_DEBUG), WithShortName(false), WithTimeStyle(FLAG_TIME_DATETIME))
}

func UseOption(logger *Logger, opts ...Option) {
	for _, opt := range opts {
		opt(logger)
	}
}

func Fatal(args ...any) {
	DEFAULT.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	DEFAULT.Fatalf(format, args...)
}

func Fatalln(args ...any) {
	DEFAULT.Fatalln(args...)
}

func Panic(args ...any) {
	DEFAULT.Panic(args...)
}

func Panicf(format string, args ...any) {
	DEFAULT.Panicf(format, args...)
}

func Panicln(args ...any) {
	DEFAULT.Panicln(args...)
}

func Print(args ...any) {
	DEFAULT.Print(args...)
}

func Printf(format string, args ...any) {
	DEFAULT.Printf(format, args...)
}

func Println(args ...any) {
	DEFAULT.Println(args...)
}

func Info(args ...any) {
	DEFAULT.Info(args...)
}

func Infof(format string, args ...any) {
	DEFAULT.Infof(format, args...)
}

func Infoln(args ...any) {
	DEFAULT.Infoln(args...)
}

func Warn(args ...any) {
	DEFAULT.Warn(args...)
}

func Warnf(format string, args ...any) {
	DEFAULT.Warnf(format, args...)
}

func Warnln(args ...any) {
	DEFAULT.Warnln(args...)
}

func Error(args ...any) {
	DEFAULT.Error(args...)
}

func Errorf(format string, args ...any) {
	DEFAULT.Errorf(format, args...)
}

func Errorln(args ...any) {
	DEFAULT.Errorln(args...)
}

func Debug(args ...any) {
	DEFAULT.Debug(args...)
}

func Debugf(format string, args ...any) {
	DEFAULT.Debugf(format, args...)
}

func Debugln(args ...any) {
	DEFAULT.Debugln(args...)
}
