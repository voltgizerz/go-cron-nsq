package main

import (
	"time"

	"github.com/voltgizerz/go-cron-nsq/config"
)

func main() {
	config.LoadENV()

	go config.InitCRON()
	go config.Consumer()
	go config.Producer()

	time.Sleep(time.Minute * 5)
}
