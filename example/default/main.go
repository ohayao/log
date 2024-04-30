package main

import (
	"math/rand"
	"time"

	"github.com/ohayao/log"
)

func main() {
	defaults()
}

func defaults() {
	log.Info("Hello world")
	log.Warn("Hello world")
	log.Error("Hello world")
	log.Debug("Hello world")
	log.Stack("Hello world")
	log.Json(log.LevelError, map[string]interface{}{"a": 1, "b": "c"}, "Hello world ")
	log.Json(log.LevelStack, map[string]interface{}{"a": 1, "b": "c"}, "Hello world ")
	log.Println("Hello world")
	if rnd()%2 == 0 {
		log.Fatalln("Hello world")
	} else {
		log.Panicln("Hello world")
	}
}

func rnd() int {
	rand.NewSource(time.Now().Unix())
	rnd := rand.Intn(100)
	return rnd
}
