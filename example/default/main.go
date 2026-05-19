package main

import (
	"fmt"
	"time"

	"github.com/ohayao/log/v2"
)

func main() {
	testDefault()
	testTerminal()
}

func testDefault() {
	log.UseOption(log.DEFAULT,
		log.WithColor(true),
		log.WithShortName(true),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_DEBUG),
	)

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Errorf("recover %v", err))
		}
	}()

	log.Info("hello world")
	log.Warnf("hello world, it's %d", time.Now().Unix())
	log.Debugln("hello world")
	log.Error("hello world")

	if time.Now().Unix()%2 == 0 {
		log.Panic("hello world")
	} else {
		log.Fatal("fatal")
	}
}

func testTerminal() {
	logger := log.NewTerminalLogger(nil,
		log.WithColor(true),
		log.WithShortName(false),
		log.WithTimeStyle(log.FLAG_TIME_TIMESTAMP),
		log.WithMinLevel(log.LV_DEBUG),
	)

	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("recover %v", err))
		}
	}()

	logger.Info("hello world")
	logger.Warnf("hello world, it's %d", time.Now().Unix())
	logger.Debugln("hello world")
	logger.Error("hello world")

	if time.Now().Unix()%2 == 0 {
		logger.Panic("hello world")
	} else {
		logger.Fatal("fatal")
	}
}
