package main

import (
	"flag"
	"log"
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	var amqpuri string

	flag.StringVar(&amqpuri, "amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpConsumer := simpleamqp.NewAmqpConsumer(amqpuri)
	messages := amqpConsumer.Receive("events", []string{"efa1", "efa2"}, "efa", 30*time.Second)
	for message := range messages {
		log.Println(message)
	}
}
