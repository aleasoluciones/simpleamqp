// +build integration

package simpleamqp_test

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

func TestPublishWithTTL(t *testing.T) {
	amqpUrl := amqpUrlFromEnv()
	amqpPublisher := NewAmqpPublisher(amqpUrl, "events")
	amqpConsumer := NewAmqpConsumer(amqpUrl)

	_, ch := Setup(amqpUrl)
	q := QueueDeclare(ch, "message_with_ttl_queue", QueueOptions{Durable: false, Delete: true, Exclusive: false})
	_ = ch.QueueBind(q.Name, "routingkey2", "events", false, nil)

	amqpPublisher.PublishWithTTL("routingkey2", []byte("irrelevantBody1"), 500)

	time.Sleep(500 * time.Millisecond)

	messages := amqpConsumer.Receive(
		"events", []string{"routingkey2"},
		"message_with_ttl_queue", QueueOptions{Durable: false, Delete: true, Exclusive: false},
		30*time.Second)

	select {
	case _ = <-messages:
		t.Error("Should not receive any message")
	case <-time.After(500 * time.Millisecond):
	}
}
