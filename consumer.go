package simpleamqp

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPConsumer represents an AMQP consumer. Used to receive messages with or without timeout
type AMQPConsumer interface {
	Receive(exchange string, routingKeys []string, queue string, queueOptions QueueOptions, queueTimeout time.Duration) chan AmqpMessage
	ReceiveWithoutTimeout(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) chan AmqpMessage
}

// AmqpConsumer holds the brokerURI
type AmqpConsumer struct {
	brokerURI string
}

// NewAmqpConsumer returns an AMQP Consumer
func NewAmqpConsumer(brokerURI string) *AmqpConsumer {
	return &AmqpConsumer{
		brokerURI: brokerURI,
	}
}

// AmqpMessage struct
type AmqpMessage struct {
	Exchange   string
	RoutingKey string
	Body       string
}

// Receive Return a AmqpMessage channel to receive messages using a given queue connected to the exchange with one ore more routing keys
// Autoreconnect on error or when we have no message after queueTimeout expired. Use 0 when not timeout is required.
// The function declares the queue
func (client *AmqpConsumer) Receive(exchange string, routingKeys []string, queue string, queueOptions QueueOptions, queueTimeout time.Duration) chan AmqpMessage {
	output := make(chan AmqpMessage)

	go func() {
		for {
			conn, ch, qname, err := client.setupConsuming(exchange, routingKeys, queue, queueOptions)
			if err != nil {
				log.Println("[simpleamqp] Error doing setupConsuming -> ", err)
				log.Println("[simpleamqp] Waiting before reconnect")
				time.Sleep(timeToReconnect)
				continue
			}
			messages, err := ch.Consume(qname, "", true, false, false, false, nil)
			if err != nil {
				log.Println("[simpleamqp] Error consuming messages -> ", err)
			}

			for closed := false; closed != true; {
				closed = messageToOuput(messages, output, queueTimeout)
			}

			log.Println("[simpleamqp] Closing connection ...")
			ch.Close()
			conn.Close()

			log.Println("[simpleamqp] Waiting before reconnect")
			time.Sleep(timeToReconnect)
		}
	}()

	return output
}

// ReceiveWithoutTimeout the same behavior that Receive method, but without using a timeout for receiving from the queue
func (client *AmqpConsumer) ReceiveWithoutTimeout(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) chan AmqpMessage {
	return client.Receive(exchange, routingKeys, queue, queueOptions, 0*time.Second)
}

func (client *AmqpConsumer) setupConsuming(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) (*amqp.Connection, *amqp.Channel, string, error) {
	conn, ch, err := setup(client.brokerURI)
	if err != nil {
		return nil, nil, "", err
	}

	err = exchangeDeclare(ch, exchange)
	if err != nil {
		return nil, nil, "", err
	}

	q, err := queueDeclare(ch, queue, queueOptions)
	if err != nil {
		return nil, nil, "", err
	}

	for _, routingKey := range routingKeys {
		err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
		if err != nil {
			log.Println("[simpleamqp] Error binding queue with the exchange")
			return nil, nil, "", err
		}
	}
	return conn, ch, q.Name, nil
}

func messageToOuput(messages <-chan amqp.Delivery, output chan AmqpMessage, queueTimeout time.Duration) (closed bool) {

	if queueTimeout == 0*time.Second {
		message, more := <-messages
		if more {
			output <- AmqpMessage{Exchange: message.Exchange, RoutingKey: message.RoutingKey, Body: string(decompress(message.Body, message.Headers))}
			return false
		}
		log.Println("[simpleamqp] No more messages... closing channel to reconnect with timeout zero")
		return true
	}

	timeoutTimer := time.NewTimer(queueTimeout)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	select {
	case message, more := <-messages:
		if more {
			output <- AmqpMessage{Exchange: message.Exchange, RoutingKey: message.RoutingKey, Body: string(decompress(message.Body, message.Headers))}
			return false
		}
		log.Println("[simpleamqp] No more messages... closing channel to reconnect")
		return true
	case <-afterTimeout:
		log.Println("[simpleamqp] Too much time without messages... closing channel to reconnect")
		return true
	}

}

func decompress(body []byte, headers map[string]interface{}) []byte {
	if headers[COMPRESS_HEADER] == true {
		gzipReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			log.Println("[simpleamqp] Error decompressing message, returning original body")
			return body
		}
		decompressedBody, _ := io.ReadAll(gzipReader)
		return decompressedBody
	}
	return body
}
