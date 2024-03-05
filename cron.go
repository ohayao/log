package log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronHandler struct {
	fd       *os.File
	fileName string
	spec     string
	lock     sync.Mutex
}

func (that *CronHandler) Write(b []byte) (n int, err error) {
	that.lock.Lock()
	defer that.lock.Unlock()
	n, err = that.fd.Write(b)
	return
}

func (that *CronHandler) Close() error {
	return that.fd.Close()
}

func (that *CronHandler) cron() {
	that.lock.Lock()
	defer that.lock.Unlock()
	// backup
	that.fd.Close()
	os.Rename(that.fileName, fmt.Sprintf("%s.bak.%s", that.fileName, time.Now().Format("060102150405")))
	that.fd, _ = os.OpenFile(that.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
}

// NewCronHandler
// file 文件位置
// spec cron标准表达式
func NewCronHandler(file string, spec string) (*CronHandler, error) {
	dir := path.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	ch := &CronHandler{
		fd:       f,
		fileName: file,
		spec:     spec,
	}
	loc, _ := time.LoadLocation("Local")
	timer := cron.New(cron.WithLocation(loc))
	_, err = timer.AddFunc(ch.spec, ch.cron)
	if err != nil {
		return nil, err
	}
	timer.Start()
	return ch, nil
}
