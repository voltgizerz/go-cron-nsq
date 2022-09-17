package config

import (
	"os"

	"github.com/jaswdr/faker"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	Cron *cron.Cron
}

func NewCRON() Cron {
	cron := cron.New()
	cron.AddFunc(os.Getenv("CRON_RULE"), cronJOB)

	return Cron{
		Cron: cron,
	}
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}

func (c *Cron) Start() {
	c.Cron.Start()
}

func cronJOB() {
	faker := faker.New()

	logger.Printf("Hi %s, This Message from CRON!", faker.Person().Name())
}
