package config

import (
	"os"

	"github.com/jaswdr/faker"
	"github.com/robfig/cron/v3"
)

func InitCRON() {
	cron := cron.New()
	cron.AddFunc(os.Getenv("CRON_RULE"), job)
	cron.Start()
}

func job() {
	faker := faker.New()

	logger.Printf("Hi %s, This Message from CRON!", faker.Person().Name())
}
