package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/voltgizerz/go-cron-nsq/config"
)

var (
	shutdown    = make(chan os.Signal, 1)
	serverError = make(chan error, 1)
	log         = config.SetupLog()
)

func main() {
	config.LoadENV()

	maxWorker := 4
	var wg sync.WaitGroup
	wg.Add(maxWorker)

	go func() {
		defer wg.Done()

		config.InitCRON()
	}()

	go func() {
		defer wg.Done()

		config.Consumer()
	}()

	go func() {
		defer wg.Done()

		config.Producer()
	}()

	go func() {
		defer wg.Done()

		terminateSignal()
	}()

	log.Println("Server Started!")

	wg.Wait() // program will wait here until all worker goroutines have reported that they're done
}

func terminateSignal() {
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	select {
	case <-shutdown:
		log.Warn("terminate signal received!")
		os.Exit(0)
	case err := <-serverError:
		log.Errorln("server error, unable to start: %v", err)
	}
}
