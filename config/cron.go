package config

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCRON() {
	cron := cron.New()
	cron.AddFunc(os.Getenv("CRON_RULE"), job)
	cron.Start()

	time.Sleep(time.Minute * 5)
}

func job() {
	log := SetupLog()

	yourName := "Felix"
	log.Println("Hi Every second! " + yourName)
}
