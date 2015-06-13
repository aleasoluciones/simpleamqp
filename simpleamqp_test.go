package simpleamqp_test

// +build integration

import (
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
	management.QueueDelete("q_initial_queueinfo")
	management.QueueDeclare("q_initial_queueinfo", QueueOptions{Durable: false, Delete: true, Exclusive: false})

	result, _ := management.QueueInfo("q_initial_queueinfo")

	assert.Equal(t, "q_initial_queueinfo", result.Name)
	assert.Equal(t, 0, result.Messages)
	assert.Equal(t, 0, result.Consumers)
}

func TestAmqpManagementCountPendingMessages(t *testing.T) {
	t.Parallel()
	amqpUrl := amqpUrlFromEnv()

	management := NewAmqpManagement(amqpUrl)
	management.QueueDelete("q_count_pending_messages")
	management.QueueDeclare("q_count_pending_messages", QueueOptions{Durable: false, Delete: true, Exclusive: false})

	amqpPublisher := NewAmqpPublisher(amqpUrl, "e1")

	management.QueueBind("q_count_pending_messages", "e1", "#")
	amqpPublisher.Publish("routingkey1", []byte("irrelevantBody1"))
	amqpPublisher.Publish("routingkey1", []byte("irrelevantBody2"))

	result, _ := management.QueueInfo("q_count_pending_messages")

	assert.Equal(t, "q_count_pending_messages", result.Name)
	assert.Equal(t, 2, result.Messages)
	assert.Equal(t, 0, result.Consumers)

}

func TestAmqpManagementCountConsumers(t *testing.T) {
	t.Parallel()
	amqpUrl := amqpUrlFromEnv()

	management := NewAmqpManagement(amqpUrl)
	management.QueueDelete("q_count_consumers")

	NewAmqpConsumer(amqpUrl).Receive("ex", []string{"#"}, "q_count_consumers",
		QueueOptions{Durable: true, Delete: false, Exclusive: false},
		30*time.Second)

	result, _ := management.QueueInfo("q_count_consumers")

	assert.Equal(t, "q_count_consumers", result.Name)
	assert.Equal(t, 0, result.Messages)
	assert.Equal(t, 1, result.Consumers)

}
