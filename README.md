# simpleamqp

[![Build Status](https://travis-ci.com/aleasoluciones/simpleamqp.svg?branch=master)](https://travis-ci.com/aleasoluciones/simpleamqp)
[![GoDoc](https://godoc.org/github.com/aleasoluciones/simpleamqp?status.png)](http://godoc.org/github.com/aleasoluciones/simpleamqp)
[![License](https://img.shields.io/github/license/aleasoluciones/http2amqp)](https://github.com/aleasoluciones/http2amqp/blob/master/LICENSE)

**simpleamqp** is a very opinionated minimal wrapper around the [AMQP](github.com/streadway/amqp) Go library with the following features:
- Publish events to a exchange 
- Consume events from a exchange
- Try to be always running
- Best effort when sending or receiving the message

## Build

You need a Go runtime installed in your system which supports [modules](https://tip.golang.org/doc/go1.16#modules). A nice way to have multiple Go versions and switch easily between them is the [g](https://github.com/stefanmaric/g) application.

A [Makefile](Makefile) is available, so you only need to run:

```sh
make build
```

Two binaries will be built, a publisher and a consumer.

## Running tests

Load environment variables to set the BROKER_URI environment variable.

```sh
source dev/env_develop
```

Start a RabbitMQ service with default configuration (specified in [`/dev/env_develop`](/dev/env_develop)).

```sh
make start_dependencies
```

Run tests. They will only work if the RabbitMQ container is up.

```sh
make test
```

## Execution example

Make sure that the environment variables are loaded before executing each command, so they can read the BROKER_URI. Also that you have a RabbitMQ container up and running with those parameters.

In one terminal, execute:

```sh
./consumer
```

In another one, execute:

```sh
./publisher
```

Now you should see how an event is sent periodically by the publisher, and read by the consumer.

## Additional info

### Publish to a exchange

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
