package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	TIME_TO_RECONNECT = 5 * time.Second
)

func setup(url string) (*amqp.Connection, *amqp.Channel) {
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Println("Error dial", err)
			time.Sleep(4 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Println("Error channel", err)
			time.Sleep(4 * time.Second)
			continue
		}
		return conn, ch
	}
}
