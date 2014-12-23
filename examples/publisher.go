package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	var amqpuri string

	flag.StringVar(&amqpuri, "amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpPublisher := simpleamqp.NewAmqpPublisher(amqpuri)
	cont := 0
	for {

		messageBody := fmt.Sprint("EFA1 ", cont)
		log.Println(messageBody)
		amqpPublisher.Publish("events", "efa1", []byte(messageBody))
		messageBody = fmt.Sprint("EFA2 ", cont)
		amqpPublisher.Publish("events", "efa2", []byte(messageBody))

		//rand.Intn(200000)
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		cont = cont + 1
	}
}
