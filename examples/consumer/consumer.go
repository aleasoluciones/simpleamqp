package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	var brokerURI string
	brokerURIParam := "brokeruri"
	brokerURIEnv := "BROKER_URI"

	flag.StringVar(&brokerURI, brokerURIParam, os.Getenv(brokerURIEnv), "Broker URI")
	flag.Parse()

	if len(brokerURI) == 0 {
		msg := fmt.Sprintf("Please specify a broker URI with the parameter -%s or set the %s environment variable.", brokerURIParam, brokerURIEnv)
		fmt.Println(msg)
		return
	}

	amqpConsumer := simpleamqp.NewAmqpConsumer(brokerURI)
	messages := amqpConsumer.ReceiveWithoutTimeout(
		"events",
		[]string{"routing_key_1", "routing_key_2"},
		"",
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
	)

	for message := range messages {
		log.Println(message)
	}
}
