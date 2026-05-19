package main

import (
	"flag"
	"time"

	"github.com/ohayao/log/v2"
)

func main() {
	dir := flag.String("d", "./logs", "dir")
	hourInterval := flag.Int("h", 1, "hours interval")
	maxAgeHours := flag.Int("m", 5, "max age hours")
	flag.Parse()

	opts := []log.Option{
		log.WithColor(false),
		log.WithShortName(false),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_INFO),
	}

	logger, _ := log.NewFileRotateLogger(*dir, *hourInterval, *maxAgeHours, opts...)

	for {
		time.Sleep(time.Second * 1)
		logger.Infof("it's now %s", time.Now().Format("06/01/02 15:04:05.000"))
		logger.Errorf("it's now %s", time.Now().Format("06/01/02 15:04:05.000"))
	}
}
