package main

import (
	"time"

	"github.com/ohayao/log/v2"
)

func main() {
	opts := []log.Option{
		log.WithColor(false),
		log.WithShortName(false),
		log.WithTimeStyle(log.FLAG_TIME_TIMESTAMP),
		log.WithMinLevel(log.LV_INFO),
	}
	// 2M
	logger, _ := log.NewFileLogger("./output/log.log", 1000*1024*2, opts...)

	for {
		time.Sleep(time.Millisecond * 10)
		logger.Infof("now is %s", time.Now().Format("06/01/02 15:04:05.000"))
		logger.Errorf("now is %s", time.Now().Format("06/01/02 15:04:05.000"))
	}
}
