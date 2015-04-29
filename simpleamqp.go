package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	TIME_TO_RECONNECT = 5 * time.Second
)

type QueueOptions struct {
	Durable   bool
	Delete    bool
	Exclusive bool
}

var DefaultQueueOptions = QueueOptions{
	Durable:   true,
	Delete:    false,
	Exclusive: false,
}

type AmqpManagement struct {
	brokerURI string
}

func NewAmqpManagement(brokerURI string) AmqpManagement {
	return AmqpManagement{brokerURI}
}

type AmqpQueueInfo struct {
	Name      string
	Messages  int
	Consumers int
}

func (i AmqpManagement) QueueDeclare(queue string, queueOptions QueueOptions) {
	conn, ch := setup(i.brokerURI)
	defer ch.Close()
	defer conn.Close()
	queueDeclare(ch, queue, queueOptions)

}

func (i AmqpManagement) QueueInfo(queue string) (AmqpQueueInfo, error) {
	conn, ch := setup(i.brokerURI)
	defer ch.Close()
	defer conn.Close()

	result, err := ch.QueueInspect(queue)
	if err != nil {
		return AmqpQueueInfo{}, err
	}
	log.Println(result)

	return AmqpQueueInfo{queue, result.Messages, result.Consumers}, nil
}

func setup(url string) (*amqp.Connection, *amqp.Channel) {
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Println("Error dial", err)
			time.Sleep(TIME_TO_RECONNECT)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Println("Error channel", err)
			time.Sleep(TIME_TO_RECONNECT)
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
			time.Sleep(TIME_TO_RECONNECT)
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
			time.Sleep(TIME_TO_RECONNECT)
			continue
		} else {
			return q
		}
	}
}
