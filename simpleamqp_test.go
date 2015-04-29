package simpleamqp_test

// +build integration

import (
	"log"
	"os"
	"time"

	"testing"

	. "github.com/aleasoluciones/simpleamqp"
	"github.com/stretchr/testify/assert"
)

func amqpUrlFromEnv() string {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://"
	}
	return url
}

func TestPublishAndReceiveTwoMessages(t *testing.T) {
	t.Parallel()
	amqpUrl := amqpUrlFromEnv()
	amqpPublisher := NewAmqpPublisher(amqpUrl, "events")
	amqpConsumer := NewAmqpConsumer(amqpUrl)
	messages := amqpConsumer.Receive(
		"events", []string{"routingkey1"},
		"", QueueOptions{Durable: false, Delete: true, Exclusive: true},
		30*time.Second)
	log.Println(amqpConsumer)

	time.Sleep(2 * time.Second)

	amqpPublisher.Publish("routingkey1", []byte("irrelevantBody1"))
	amqpPublisher.Publish("routingkey1", []byte("irrelevantBody2"))

	message1 := <-messages
	assert.Equal(t, message1.Body, "irrelevantBody1")
	assert.Equal(t, message1.Exchange, "events")
	assert.Equal(t, message1.RoutingKey, "routingkey1")
	message2 := <-messages
	assert.Equal(t, message2.Body, "irrelevantBody2")
	assert.Equal(t, message2.Exchange, "events")
	assert.Equal(t, message2.RoutingKey, "routingkey1")

}

func TestAmqpManagementInitialQueueInfo(t *testing.T) {
	t.Parallel()
	amqpUrl := amqpUrlFromEnv()

	management := NewAmqpManagement(amqpUrl)
	management.QueueDeclare("events", QueueOptions{Durable: false, Delete: true, Exclusive: false})

	result, _ := management.QueueInfo("events")

	assert.Equal(t, result.Name, "events")
	assert.Equal(t, result.Messages, 0)
	assert.Equal(t, result.Consumers, 0)
}
