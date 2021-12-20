package simpleamqp

import (
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

const (
	timeToWaitForChannel = 5 * time.Second // seconds to wait and don't block when writting to channel publisher.outputMessages
)

type messageToPublish struct {
	routingKey string
	message    []byte
	expiration string
}

// AMQPPublisher represents an AMQP Publisher that can publish messages with or without TTL
type AMQPPublisher interface {
	Publish(string, []byte)
	PublishWithTTL(string, []byte, int)
}

// AmqpPublisher holds the brokerURI, exchange name and channel where to submit messages to be publish to rabbitmq
type AmqpPublisher struct {
	brokerURI      string
	exchange       string
	outputMessages chan messageToPublish
}

// NewAmqpPublisher returns an AmqpPublisher
func NewAmqpPublisher(brokerURI, exchange string) *AmqpPublisher {
	publisher := AmqpPublisher{
		brokerURI:      brokerURI,
		exchange:       exchange,
		outputMessages: make(chan messageToPublish, 1024),
	}

	go func() {
		for {
			err := publisher.publishLoop()
			log.Println("[simpleamqp] Waiting", timeToReconnect, "to reconnect due ", err)
			time.Sleep(timeToReconnect)
		}
	}()
	return &publisher
}

// Publish publish a message using the given routing key
func (publisher *AmqpPublisher) Publish(routingKey string, message []byte) {
	publisher.queueMessageToPublish(messageToPublish{routingKey: routingKey, message: message})
}

// PublishWithTTL publish a message waiting the given TTL
func (publisher *AmqpPublisher) PublishWithTTL(routingKey string, message []byte, ttl int) {
	publisher.queueMessageToPublish(messageToPublish{routingKey: routingKey, message: message, expiration: strconv.Itoa(ttl)})
}

// Queue the message to be published and return inmediatly
// The message will be published to the AmqpPublisher exchange using the given routingKey
// If the message can't be queued after some short time (because the channel is full) a log is printed and the message is discarded
func (publisher *AmqpPublisher) queueMessageToPublish(messageToPublish messageToPublish) {
	timeoutTimer := time.NewTimer(timeToWaitForChannel)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	select {
	case publisher.outputMessages <- messageToPublish:
	case <-afterTimeout:
		log.Println("[simpleamqp] Publish channel full", messageToPublish)
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
			Expiration:      messageToPublish.expiration,
		})
	return err

}

func (publisher *AmqpPublisher) publishLoop() error {
	conn, ch, err := setup(publisher.brokerURI)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()

	err = exchangeDeclare(ch, publisher.exchange)
	if err != nil {
		return err
	}
	for {
		messageToPublish := <-publisher.outputMessages
		err := publisher.publish(ch, messageToPublish)
		if err != nil {
			return err
		}
		log.Println("[simpleamqp] Published", messageToPublish.routingKey, string(messageToPublish.message))
	}
}
