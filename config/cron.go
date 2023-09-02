package config

import (
	"github.com/robfig/cron/v3"
)

type Cron struct {
	Cron     *cron.Cron
	Producer *Producer
}

func NewCRON(p *Producer) Cron {
	cron := cron.New()

	return Cron{
		Cron:     cron,
		Producer: p,
	}
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}

func (c *Cron) Start() {
	c.Cron.Start()
}

func (c *Cron) JOB() {
	// faker := faker.New()

	// logger.Printf("Hi %s, This Message from CRON!", faker.Person().Name())
	c.Producer.Publish("topic")
}
