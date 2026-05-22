package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type fileRotateHandler struct {
	lock           sync.Mutex
	fd             *os.File
	file           string
	maxAgeHours    int       // 最大存储小时
	hoursInterval  int       // 每几小时
	lastRotateTime time.Time // 上次轮转时间
}

func (f *fileRotateHandler) Write(b []byte) (n int, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.check()
	n, err = f.fd.Write(b)
	return
}

func (f *fileRotateHandler) Close() (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.fd.Close()
}

func (f *fileRotateHandler) check() {
	now := time.Now()
	if f.lastRotateTime.IsZero() {
		f.lastRotateTime = now.Truncate(time.Hour)
		return
	}

	nextRotateTime := f.lastRotateTime.Add(time.Duration(f.hoursInterval) * time.Hour)
	if now.Before(nextRotateTime) {
		return
	}
	dir, fileName := GetDirAndFileName(f.file, "log.log")
	bakFileName := fmt.Sprintf("%sbak_%s_%s", dir, f.lastRotateTime.Format("2006010215"), fileName)
	_ = f.fd.Close()
	_ = os.Rename(f.file, bakFileName)
	f.fd, _ = os.OpenFile(f.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	f.lastRotateTime = nextRotateTime

	go f.cleanOldFiles(dir, now)
}

func (f *fileRotateHandler) cleanOldFiles(dir string, now time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	maxAgeTime := now.Add(-time.Duration(f.maxAgeHours) * time.Hour)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), "bak_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(maxAgeTime) {
			_ = os.Remove(filepath.Join(dir, entry.Name()))
		}
	}
}

func newFileRotateHandler(file string, hoursInterval, maxAgeHours int) (*fileRotateHandler, error) {
	dir, _file := GetDirAndFileName(file, "log.log")
	handler := &fileRotateHandler{
		file:          dir + _file,
		hoursInterval: hoursInterval,
		maxAgeHours:   maxAgeHours,
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(handler.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	handler.fd = f
	return handler, nil
}
