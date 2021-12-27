package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

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

	amqpPublisher := simpleamqp.NewAmqpPublisher(brokerURI, "events")
	cont := 0
	for {
		payload := fmt.Sprintf("Message #%d", cont)
		amqpPublisher.Publish("routing_key_1", []byte(payload))
		amqpPublisher.Publish("routing_key_2", []byte(payload))

		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		cont = cont + 1
	}
}
