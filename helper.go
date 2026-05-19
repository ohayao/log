package log

import (
	"path/filepath"
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

func ShortFileName(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx > -1 {
		return path[idx+1:]
	}
	return path
}

// GetDirAndFileName 获取目录和文件名
//
// dir 目录，尾部带"/"
//
// fileName 文件名，如果路径中没有文件名，则使用defaultFileName
func GetDirAndFileName(path string, defaultFileName string) (dir, fileName string) {
	path = strings.ReplaceAll(path, `\`, "/")
	if path == "" || path == "." || path == "./" {
		return "./", defaultFileName
	}
	isDir := strings.HasSuffix(path, "/")
	path = strings.TrimRight(path, "/")
	if path == "" {
		return "/", defaultFileName
	}
	if isDir || path == ".." {
		return path + "/", defaultFileName
	}
	dir = filepath.Dir(path)
	fileName = filepath.Base(path)
	dir = strings.TrimRight(dir, "/") + "/"
	if fileName == "" || fileName == "." {
		fileName = defaultFileName
	}
	return
}
