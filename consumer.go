package simpleamqp

import (
	"log"
	"time"
)

type AMQPConsumer interface {
	Receive(exchange string, routingKeys []string, queue string, queueTimeout time.Duration) chan AmqpMessage
}

type AmqpConsumer struct {
	brokerUri string
}

// Return AMQP Consumer
func NewAmqpConsumer(brokerUri string) *AmqpConsumer {
	return &AmqpConsumer{
		brokerUri: brokerUri,
	}
}

// AmqpMessage struct
type AmqpMessage struct {
	Body string
}

// Return a AmqpMessage channel to receive messages using a given queue connected to the exchange with one ore more routing keys
// Autoreconnect on error or when we have no message after queueTimeout expired
// The function declare the queue
func (client *AmqpConsumer) Receive(exchange string, routingKeys []string, queue string, queueTimeout time.Duration) chan AmqpMessage {
	output := make(chan AmqpMessage)

	go func() {
		for {
			conn, ch := setup(client.brokerUri)
			defer conn.Close()
			defer ch.Close()

			exchangeDeclare(ch, exchange)
			q := queueDeclare(ch, queue)

			for _, routingKey := range routingKeys {
				_ = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
			}

			messages, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

			for closed := false; closed != true; {
				select {
				case message, more := <-messages:
					if more {
						output <- AmqpMessage{Body: string(message.Body)}
					} else {
						log.Println("[simpleamqp] No more messages... closing channel to reconnect")
						closed = true
					}
				case <-time.After(queueTimeout):
					log.Println("[simpleamqp] Too much time without messages... closing channel to reconnect")
					closed = true
				}
			}
			log.Println("[simpleamqp] Waiting befor reconnect")
			time.Sleep(TIME_TO_RECONNECT)
		}
	}()

	return output
}
