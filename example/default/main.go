package main

import (
	"fmt"
	"time"

	"github.com/ohayao/log/v2"
)

func main() {

	testTerminal()
	testFile()
	testRotateFile()

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Errorf("recover %v", err))
		}
	}()

	if time.Now().Unix()%2 == 0 {
		log.Panic("hello world")
	} else {
		log.Fatal("fatal")
	}
}

func testTerminal() {
	logger := log.NewTerminalLogger(nil,
		log.WithColor(true),
		log.WithShortName(true),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_DEBUG),
	)
	logger.Info("hello world")
	logger.Warnf("hello world, it's %d", time.Now().Unix())
	logger.Debugln("hello world")
	logger.Error("hello world")
}

func testFile() {
	logger, _ := log.NewFileLogger("./log/log.log", 1024*60,
		log.WithShortName(false),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_INFO))
	logger.Info("hello world")
	logger.Warnf("hello world, it's %d", time.Now().Unix())
	logger.Debugln("hello world")
	logger.Error("hello world")
}

func testRotateFile() {
	logger, _ := log.NewFileRotateLogger("./log2", "log.log", 24*3, 1,
		log.WithShortName(false),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_INFO))
	logger.Info("hello world")
	logger.Warnf("hello world, it's %d", time.Now().Unix())
	logger.Debugln("hello world")
	logger.Error("hello world")
}
