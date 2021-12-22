//go:build integration

package simpleamqp

import (
	"os"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func amqpURLFromEnv() string {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://"
	}
	return url
}

func TestPublishAndReceiveTwoMessages(t *testing.T) {
	amqpURL := amqpURLFromEnv()
	amqpPublisher := NewAmqpPublisher(amqpURL, "events")
	amqpConsumer := NewAmqpConsumer(amqpURL)
	messages := amqpConsumer.Receive(
		"events", []string{"routingkey1"},
		"", QueueOptions{Durable: false, Delete: true, Exclusive: true},
		30*time.Second)

	// Sleep sometime so the consumer can create the queue and bind.
	// Then why in the TTL test this sleep is not needed? Because
	// we are creating the queue in a synchronous way using the private functions.
	time.Sleep(7 * time.Second)

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

func TestPublishWithTTL(t *testing.T) {
	amqpURL := amqpURLFromEnv()
	exchange := "events"
	queueName := "message_with_ttl_queue"
	queueOptions := QueueOptions{Durable: false, Delete: true, Exclusive: false}
	_, ch, err := setup(amqpURL)
	if err != nil {
		t.Error("There was an error getting the channel")
	}
	_, err = queueDeclare(ch, queueName, queueOptions)
	if err != nil {
		t.Error("There was an error declaring the queue")
	}
	err = ch.QueueBind(queueName, "routingkey2", exchange, false, nil)
	if err != nil {
		t.Error("There was an error binding the exchange with the queue using the given routing key")
	}

	messageTTL := 500
	amqpPublisher := NewAmqpPublisher(amqpURL, exchange)
	amqpPublisher.PublishWithTTL("routingkey2", []byte("irrelevantBody1"), messageTTL)

	time.Sleep((time.Duration(messageTTL + 13)) * time.Millisecond)

	amqpConsumer := NewAmqpConsumer(amqpURL)
	messages := amqpConsumer.Receive(
		exchange, []string{"routingkey2"},
		queueName, queueOptions,
		30*time.Second)

	select {
	case <-messages:
		t.Error("Should not receive any message")
	case <-time.After(500 * time.Millisecond):
	}
}
