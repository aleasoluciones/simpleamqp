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

func setup(url string) (*amqp.Connection, *amqp.Channel) {
	for {
		conn, err := amqp.Dial(normalizeURL(url))
		if err != nil {
			log.Println("Error dial", err)
			time.Sleep(timeToReconnect)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Println("Error channel", err)
			conn.Close()
			time.Sleep(timeToReconnect)
			continue
		}
		return conn, ch
	}
}

func exchangeDeclare(ch *amqp.Channel, exchange string) {
	for {
		log.Println("Exchange declare", exchange)
		err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
		if err != nil {
			log.Println("Error declaring exchange", err)
			time.Sleep(timeToReconnect)
			continue
		} else {
			return
		}
	}
}

func queueDeclare(ch *amqp.Channel, queue string, queueOptions QueueOptions) amqp.Queue {
	for {
		log.Println("Queue declare", queue)
		q, err := ch.QueueDeclare(
			queue, // name of the queue
			queueOptions.Durable,
			queueOptions.Delete,
			queueOptions.Exclusive,
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			log.Println("Error declaring queue", err)
			time.Sleep(timeToReconnect)
			continue
		} else {
			return q
		}
	}
}
