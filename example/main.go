package main

import (
	"os"
	"time"

	"github.com/ohayao/log"
)

var logger *log.Logger

func init() {
	handler := log.NewStreamHandler(os.Stderr)
	logger = log.NewLogger(handler)
	logger.SetFlags(log.FlagTime, log.FlagLevel, log.FlagColor)
	logger.SetLevels(log.LevelAll)
}

func main() {
	defer logger.Close()
	logger.Info(time.Now())
	logger.Warn(time.Now())
	logger.Error(time.Now())
	logger.Debug(time.Now())
	logger.Stackln(log.DefaultDepth, time.Now())
	logger.Json(log.LevelError, map[string]interface{}{"a": 1, "b": "c"}, "json ")
	logger.Println(time.Now())
	if time.Now().Second()%2 == 0 {
		logger.Fatalln(time.Now())
	} else {
		logger.Panicln(time.Now())
	}
}
