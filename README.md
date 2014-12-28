#simpleamqp
[![Build Status](https://travis-ci.org/aleasoluciones/simpleamqp.svg?branch=master)](https://travis-ci.org/aleasoluciones/simpleamqp)
[![GoDoc](https://godoc.org/github.com/aleasoluciones/simpleamqp?status.png)](http://godoc.org/github.com/aleasoluciones/simpleamqp)

SimpleAMQP is a minimal wrapper around the excelent AMQP library [github.com/streadway/amqp](http://github.com/streadway/amqp)
It provided a AMQP Consumer and a AMQP Publisher
## Uses Cases
### Publish to a exchange
Publish messages to a exchange without blocking the producing
**Features**
 * Reconnect when a messege can't be delvered
 * Message buffer to avoid blocking the publisher
 * Discard messages when the message buffer is full due to a connection problem



## Tests
To run the integration tests, make sure you have RabbitMQ running on any host, export the environment variable AMQP_URL=amqp://host/ and run go test -tags integration. TravisCI will also run the integration tests.