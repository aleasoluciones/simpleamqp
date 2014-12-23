package main

import (
	"flag"
	"log"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	var amqpuri string

	flag.StringVar(&amqpuri, "amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpClient := simpleamqp.NewAmqpConsumer(amqpuri)
	messages := amqpClient.Receive("events", []string{"efa1", "efa2"}, "efa")
	for message := range messages {
		log.Println(message)
	}
}
