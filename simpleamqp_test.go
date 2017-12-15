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
	exchange := "events"
	queueName := "message_with_ttl_queue"
	queueOptions := QueueOptions{Durable: false, Delete: true, Exclusive: false}
	_, ch := Setup(amqpUrl)
	QueueDeclare(ch, queueName, queueOptions)
	_ = ch.QueueBind(queueName, "routingkey2", exchange, false, nil)

	messageTTL := 500
	amqpPublisher := NewAmqpPublisher(amqpUrl, exchange)
	amqpPublisher.PublishWithTTL("routingkey2", []byte("irrelevantBody1"), messageTTL)

	time.Sleep((time.Duration(messageTTL) + 13 ) * time.Millisecond)

	amqpConsumer := NewAmqpConsumer(amqpUrl)
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
