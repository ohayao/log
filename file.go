package log

import (
	"fmt"
	"math"
	"os"
	"path"
	"time"
)

type FileHandler struct {
	fd               *os.File
	fileName         string
	maxSize, curSize int64
}

func (that *FileHandler) Write(b []byte) (n int, err error) {
	that.check()
	n, err = that.fd.Write(b)
	that.curSize += int64(n)
	return
}

func (that *FileHandler) Close() error {
	return that.fd.Close()
}

func (that *FileHandler) check() {
	if that.maxSize > that.curSize {
		return
	}
	stat, err := that.fd.Stat()
	if err != nil {
		return
	}
	if stat.Size() < that.maxSize {
		that.curSize = stat.Size()
		return
	}
	// backup
	that.fd.Close()
	os.Rename(that.fileName, fmt.Sprintf("%s.bak.%s", that.fileName, time.Now().Format("060102150405")))
	that.fd, _ = os.OpenFile(that.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	that.curSize = 0
	f, err := that.fd.Stat()
	if err != nil {
		return
	}
	that.curSize = f.Size()
}

// NewFileHandler
// file 文件位置
// maxSize 文件大小，单位Bytes
func NewFileHandler(file string, maxSize int64) (*FileHandler, error) {
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
	return &FileHandler{
		fd:       f,
		fileName: file,
		maxSize:  maxSize,
		curSize:  stat.Size(),
	}, nil
}
