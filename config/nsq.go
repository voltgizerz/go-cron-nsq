package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaswdr/faker"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	Name    string
	Message string
}

func createMessage() []byte {
	faker := faker.New()

	name := faker.Person().Name()
	message := Message{
		Name:    name,
		Message: fmt.Sprintf("Hi %s, This Message from Producer!", name),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		logger.Error("Error when marshal message body!")
	}

	return messageBody
}

type Producer struct {
	Client *nsq.Producer
}

func NewProducer(addr, topicName string) *Producer {
	// Instantiate a producer.
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, config)
	if err != nil {
		logger.Error(err)
	}

	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	for i := 0; i < 10; i++ {
		err = producer.Publish(topicName, createMessage())
		if err != nil {
			logger.Error(err)
		} else {
			logger.Printf(" [X] Successfully Published Message %d...", i+1)
		}
	}

	return &Producer{
		Client: producer,
	}
}

func (p *Producer) Publish(topicName string) {
	err := p.Client.Publish(topicName, createMessage())
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Printf(" [X] Successfully Published Message From CRON...")
}

func (p *Producer) Stop() {
	// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
	p.Client.Stop()
}

// HandleMessage implements the Handler interface.
func (h *Message) HandleMessage(m *nsq.Message) error {
	var err error
	if len(m.Body) == 0 {
		return nil
	}

	// do whatever actual message processing is desired
	logger.Printf(" [X] Received message : %s Retry attempt %d", m.Body, m.Attempts)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}

// Consumer
type Consumer struct {
	Client *nsq.Consumer
}

// NewConsumer - function.
func NewConsumer(addr, channel, topicName string) *Consumer {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topicName, channel, config)
	if err != nil {
		logger.Error(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(&Message{})

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		logger.Error(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	defer consumer.Stop()

	return &Consumer{
		Client: consumer,
	}
}

func (c *Consumer) Stop() {
	// Gracefully stop the consumer.
	c.Client.Stop()
}
