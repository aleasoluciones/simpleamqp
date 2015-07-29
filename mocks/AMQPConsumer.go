package mocks

import (
	"github.com/aleasoluciones/simpleamqp"
	"github.com/stretchr/testify/mock"
)

import "time"

type AMQPConsumer struct {
	mock.Mock
}

func (m *AMQPConsumer) Receive(exchange string, routingKeys []string, queue string, options simpleamqp.QueueOptions, queueTimeout time.Duration) chan simpleamqp.AmqpMessage {
	ret := m.Called(exchange, routingKeys, queue, queueTimeout)
	return ret.Get(0).(chan simpleamqp.AmqpMessage)
}
