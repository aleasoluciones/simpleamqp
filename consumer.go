package simpleamqp

import (
	"log"
	"time"

	"github.com/kr/pretty"
	//"github.com/streadway/amqp"
)

type AmqpConsumer struct {
	brokerUri string
}

func NewAmqpConsumer(brokerUri string) *AmqpConsumer {
	return &AmqpConsumer{
		brokerUri: brokerUri,
	}
}

type AmqpMessage struct {
	Body string
}

func (client *AmqpConsumer) Receive(exchange string, routingKeys []string, queue string) chan AmqpMessage {
	output := make(chan AmqpMessage)

	go func() {
		for {
			conn, ch := setup(client.brokerUri)
			defer conn.Close()
			defer ch.Close()

			//channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
			q, _ := ch.QueueDeclare(queue, true, false, false, false, nil)
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
						pretty.Println("No more messages... closing channel to reconnect")
						closed = true
					}
				case <-time.After(5 * time.Second):
					log.Println("Too much time without messages... closing channel to reconnect")
					closed = true
				}
			}
			log.Println("Waiting befor reconnect")
			time.Sleep(5 * time.Second)
		}
	}()

	return output
}
