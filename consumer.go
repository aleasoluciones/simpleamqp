package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type AMQPConsumer interface {
	Receive(exchange string, routingKeys []string, queue string, queueTimeout time.Duration) chan AmqpMessage
}

type AmqpConsumer struct {
	brokerURI string
}

// Return AMQP Consumer
func NewAmqpConsumer(brokerURI string) *AmqpConsumer {
	return &AmqpConsumer{
		brokerURI: brokerURI,
	}
}

// AmqpMessage struct
type AmqpMessage struct {
	Body string
}

// Return a AmqpMessage channel to receive messages using a given queue connected to the exchange with one ore more routing keys
// Autoreconnect on error or when we have no message after queueTimeout expired
// The function declares the queue
func (client *AmqpConsumer) Receive(exchange string, routingKeys []string, queue string, queueOptions QueueOptions, queueTimeout time.Duration) chan AmqpMessage {
	output := make(chan AmqpMessage)

	go func() {
		for {
			conn, ch := setup(client.brokerURI)

			exchangeDeclare(ch, exchange)
			q := queueDeclare(ch, queue, queueOptions)

			for _, routingKey := range routingKeys {
				_ = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
			}

			messages, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

			for closed := false; closed != true; {
				closed = messageToOuput(messages, output, queueTimeout)
			}

			log.Println("[simpleamqp] Closing connection ...")
			ch.Close()
			conn.Close()

			log.Println("[simpleamqp] Waiting befor reconnect")
			time.Sleep(TIME_TO_RECONNECT)
		}
	}()

	return output
}

func messageToOuput(messages <-chan amqp.Delivery, output chan AmqpMessage, queueTimeout time.Duration) (closed bool) {
	timeoutTimer := time.NewTimer(queueTimeout)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	detectedClosed := false
	select {
	case message, more := <-messages:
		if more {
			output <- AmqpMessage{Body: string(message.Body)}
		} else {
			log.Println("[simpleamqp] No more messages... closing channel to reconnect")
			detectedClosed = true
		}
	case <-afterTimeout:
		log.Println("[simpleamqp] Too much time without messages... closing channel to reconnect")
		detectedClosed = true
	}
	return detectedClosed
}
