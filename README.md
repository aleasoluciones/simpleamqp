# simpleamqp

[![Build Status](https://travis-ci.com/aleasoluciones/simpleamqp.svg?branch=master)](https://travis-ci.com/aleasoluciones/simpleamqp)
[![GoDoc](https://godoc.org/github.com/aleasoluciones/simpleamqp?status.png)](http://godoc.org/github.com/aleasoluciones/simpleamqp)

SimpleAMQP is very opinionated minimal wrapper around the excelent AMQP library [github.com/streadway/amqp](http://github.com/streadway/amqp), created to help implement microservices with the following characteristic:
 * Publish all the domain especific events to a exchange 
 * Reacts to events from a exchange
 * Long running processes that tries to be always running
 * Best effort to send or receive the message

## Uses Cases
### Publish to a exchange
Publish messages to a exchange without blocking the producing
#### Features
 * Reconnect when a message can't be delivered
 * Message buffer to avoid blocking the publisher
 * Discard messages when the message buffer is full due to a connection problem

#### Known Issues
 * When there is a connection problem, some messages can be lost

#### Unimplemented features
 * Exchange options not configurable
 * Message headers and characteristics (delivery mode, ttl, ...) not configurable

### Receive messages from a exchange
#### Features
 * Receive using various routing keys
 * main option for the queue configurables
 * Reconnect in case of connection error
 * Reconnect when no messages received in a configurable amount of time (to detect some kind of connection problems at NATED networks where the NATED connection expires)

#### Known Issues
 * When there is a connection problem, some messages can be lost

#### Unimplemented features
 * Exchange options not configurable

## Tests
To run the integration tests, make sure you have RabbitMQ running on any host, export the environment variable AMQP_URL=amqp://host/ and run **go test -tags integration**. TravisCI will also run the integration tests.
