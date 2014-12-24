package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type messageToPublish struct {
	exchange   string
	routingKey string
	message    []byte
}

type AmqpPublisher struct {
	brokerUri      string
	outputMessages chan messageToPublish
}

func NewAmqpPublisher(brokerUri string) *AmqpPublisher {
	publisher := AmqpPublisher{
		brokerUri:      brokerUri,
		outputMessages: make(chan messageToPublish, 1024),
	}

	go func() {
		for {
			err := publisher.publish_loop()
			log.Println("ERROR", err)
			log.Println("Waiting", TIME_TO_RECONNECT, "to reconnect")
			time.Sleep(TIME_TO_RECONNECT)
		}
	}()
	return &publisher
}

func (client *AmqpPublisher) Publish(exchange string, routingKey string, message []byte) {
	messageToPublish := messageToPublish{exchange, routingKey, message}
	select {
	case client.outputMessages <- messageToPublish:
	case <-time.After(5 * time.Second):
		log.Println("Publish channel full", messageToPublish)
	}
}

func publish(channel *amqp.Channel, messageToPublish messageToPublish) error {
	err := channel.Publish(
		messageToPublish.exchange,
		messageToPublish.routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            messageToPublish.message,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		})
	return err

}

func (client *AmqpPublisher) publish_loop() error {
	conn, ch := setup(client.brokerUri)
	defer conn.Close()
	defer ch.Close()
	for {
		messageToPublish := <-client.outputMessages
		err := publish(ch, messageToPublish)
		if err != nil {
			return err
		}
		log.Println("Published", string(messageToPublish.message))
	}
}
