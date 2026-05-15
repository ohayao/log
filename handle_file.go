package log

import (
	"fmt"
	"math"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type fileHandler struct {
	lock     sync.Mutex
	fd       *os.File
	fileName string
	maxSize  int64
	curSize  atomic.Int64
}

func (f *fileHandler) Write(b []byte) (n int, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.check()
	n, err = f.fd.Write(b)
	f.curSize.Add(int64(n))
	return
}

func (f *fileHandler) Close() (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.fd.Close()
}

func (f *fileHandler) check() {
	if f.maxSize > f.curSize.Load() {
		return
	}
	stat, err := f.fd.Stat()
	if err != nil {
		return
	}
	if stat.Size() < f.maxSize {
		f.curSize.Store(stat.Size())
		return
	}
	// backup
	_ = f.fd.Close()
	os.Rename(f.fileName, fmt.Sprintf("%s.bak.%s", f.fileName, time.Now().Format("060102150405")))
	f.fd, _ = os.OpenFile(f.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	f.curSize.Store(0)
	fi, err := f.fd.Stat()
	if err != nil {
		return
	}
	f.curSize.Store(fi.Size())
}

func newFileHandler(file string, maxSize int64) (*fileHandler, error) {
	dir := path.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if maxSize < 1 {
		maxSize = math.MaxInt64 - 1024
	}
	handler := &fileHandler{
		fd:       f,
		fileName: file,
		maxSize:  maxSize,
	}
	handler.curSize.Store(stat.Size())
	return handler, nil
}
