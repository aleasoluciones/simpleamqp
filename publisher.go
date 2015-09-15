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

type AMQPPublisher interface {
	Publish(routingKey string, message []byte)
}

type AmqpPublisher struct {
	brokerURI      string
	exchange       string
	outputMessages chan messageToPublish
}

func NewAmqpPublisher(brokerURI, exchange string) *AmqpPublisher {
	publisher := AmqpPublisher{
		brokerURI:      brokerURI,
		exchange:       exchange,
		outputMessages: make(chan messageToPublish, 1024),
	}

	go func() {
		for {
			err := publisher.publishLoop()
			log.Println("Error", err)
			log.Println("Waiting", timeToReconnect, "to reconnect")
			time.Sleep(timeToReconnect)
		}
	}()
	return &publisher
}

// Queue the message to be published and return inmediatly
// The message will be published to the AmqpPublisher exchange using the given routingKey
// If the message can't be queued (because the channel is full) a log is printed and the message is discarded
func (publisher *AmqpPublisher) Publish(routingKey string, message []byte) {
	messageToPublish := messageToPublish{routingKey, message}

	timeoutTimer := time.NewTimer(5 * time.Second)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	select {
	case publisher.outputMessages <- messageToPublish:
	case <-afterTimeout:
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

func (publisher *AmqpPublisher) publishLoop() error {
	conn, ch := setup(publisher.brokerURI)
	defer conn.Close()
	defer ch.Close()

	exchangeDeclare(ch, publisher.exchange)
	for {
		messageToPublish := <-publisher.outputMessages
		err := publisher.publish(ch, messageToPublish)
		if err != nil {
			return err
		}
		log.Println("Published", messageToPublish.routingKey, string(messageToPublish.message))
	}
}
