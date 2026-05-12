package main

import (
	"fmt"
	"time"

	"github.com/ohayao/log/v2"
)

func main() {
	log.Println("hello world")
	log.Info("hello world")
	log.Warnf("hello world, it's %d", time.Now().Unix())
	log.Debugln("hello world")
	log.Error("hello world")

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
