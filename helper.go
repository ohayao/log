package log

import (
	"runtime"
	"strings"
)

func WhoCalledMe() (file string, line int, fn string) {
	var pc uintptr
	var ok bool
	for skip := 1; ; skip++ {
		pc, file, line, ok = runtime.Caller(skip)
		if !ok {
			break
		}
		fnc := runtime.FuncForPC(pc)
		if fnc == nil {
			continue
		}
		fn = fnc.Name()
		if !strings.HasPrefix(fn, "github.com/ohayao/log") {
			if idx := strings.LastIndex(fn, "."); idx > -1 {
				fn = fn[idx+1:]
			}
			return
		}
	}
	return "?", 0, "?"
}

func ShortFileName(filePath string) string {
	idx := strings.LastIndex(filePath, "/")
	if idx > -1 {
		return filePath[idx+1:]
	}
	return filePath
}
