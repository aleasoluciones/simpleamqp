package simpleamqp

import (
	"log"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

const (
	timeToReconnect = 5 * time.Second
)

// QueueOptions holds the flags when declaring a queue:
// durable, delete, exclusive
type QueueOptions struct {
	Durable   bool
	Delete    bool
	Exclusive bool
}

// DefaultQueueOptions define the default options when creating queues:
// durable, not autodelete and not exclusive
var DefaultQueueOptions = QueueOptions{
	Durable:   true,
	Delete:    false,
	Exclusive: false,
}

func normalizeURL(url string) string {
	return strings.Replace(url, "rabbitmq://", "amqp://", -1)
}

func setup(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(normalizeURL(url))
	if err != nil {
		log.Println("[simpleamqp] Error dial", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("[simpleamqp] Error channel", err)
		conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}

func exchangeDeclare(ch *amqp.Channel, exchange string) error {
	log.Println("[simpleamqp] Exchange declare", exchange)
	err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		log.Println("[simpleamqp] Error declaring exchange", err)
	}
	return err
}

func queueDeclare(ch *amqp.Channel, queue string, queueOptions QueueOptions) (amqp.Queue, error) {
	log.Println("[simpleamqp] Queue declare", queue)
	q, err := ch.QueueDeclare(
		queue, // name of the queue
		queueOptions.Durable,
		queueOptions.Delete,
		queueOptions.Exclusive,
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Println("[simpleamqp] Error declaring queue", err)
	}
	return q, err
}
