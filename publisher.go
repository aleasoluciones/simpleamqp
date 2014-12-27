package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type messageToPublish struct {
	routingKey string
	message    []byte
}

type AmqpPublisher struct {
	brokerUri      string
	exchange       string
	outputMessages chan messageToPublish
}

func NewAmqpPublisher(brokerUri, exchange string) *AmqpPublisher {
	publisher := AmqpPublisher{
		brokerUri:      brokerUri,
		exchange:       exchange,
		outputMessages: make(chan messageToPublish, 1024),
	}

	go func() {
		for {
			err := publisher.publish_loop()
			log.Println("Error", err)
			log.Println("Waiting", TIME_TO_RECONNECT, "to reconnect")
			time.Sleep(TIME_TO_RECONNECT)
		}
	}()
	return &publisher
}

func (publisher *AmqpPublisher) Publish(routingKey string, message []byte) {
	messageToPublish := messageToPublish{routingKey, message}
	select {
	case publisher.outputMessages <- messageToPublish:
	case <-time.After(5 * time.Second):
		log.Println("Publish channel full", messageToPublish)
	}
}

func (publisher *AmqpPublisher) publish(channel *amqp.Channel, messageToPublish messageToPublish) error {
	err := channel.Publish(
		publisher.exchange,
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

func (publisher *AmqpPublisher) publish_loop() error {
	conn, ch := setup(publisher.brokerUri)
	defer conn.Close()
	defer ch.Close()

	log.Println("Exchange declare", publisher.exchange)
	err := ch.ExchangeDeclare(publisher.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		messageToPublish := <-publisher.outputMessages
		err := publisher.publish(ch, messageToPublish)
		if err != nil {
			return err
		}
		log.Println("Published", string(messageToPublish.message))
	}
}
