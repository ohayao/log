package log

import (
	"fmt"
	"math"
	"os"
	"path"
	"sync/atomic"
	"time"
)

type fileHandler struct {
	fd       *os.File
	fileName string
	maxSize  int64
	curSize  atomic.Int64
}

func (f *fileHandler) Write(b []byte) (n int, err error) {
	f.check()
	n, err = f.fd.Write(b)
	f.curSize.Add(int64(n))
	return
}

func (f *fileHandler) Close() (err error) {
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

type terminalHandler struct {
	w *os.File
}

func (t *terminalHandler) Write(b []byte) (n int, err error) {
	return t.w.Write(b)
}

func (t *terminalHandler) Close() (err error) {
	return t.w.Close()
}

// if file is nil, use os.Stdout
func newTerminalHandler(file *os.File) *terminalHandler {
	if file == nil {
		file = os.Stdout
	}
	return &terminalHandler{w: file}
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

func New(handler IHandler, opts ...Option) *Logger {
	logger := &Logger{
		handler:     handler,
		enableColor: true,
		shortName:   false,
		flagTime:    FLAG_TIME_DATETIME,
		level:       LV_DEBUG,
		pool:        poolNew(),
	}
	for _, opt := range opts {
		opt(logger)
	}

	switch handler.(type) {
	case *fileHandler:
		logger.enableColor = false
	}
	return logger
}

func NewTerminalLogger(file *os.File, opts ...Option) *Logger {
	return New(newTerminalHandler(file), opts...)
}

func NewFileLogger(file string, maxSize int64, opts ...Option) (*Logger, error) {
	handler, err := newFileHandler(file, maxSize)
	if err != nil {
		return nil, err
	}
	return New(handler, opts...), nil
}
