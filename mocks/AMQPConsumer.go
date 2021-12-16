// Package mocks represent the AMQPConsumer mocks
// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import simpleamqp "github.com/aleasoluciones/simpleamqp"
import time "time"

// AMQPConsumer is an autogenerated mock type for the AMQPConsumer type
type AMQPConsumer struct {
	mock.Mock
}

// Receive provides a mock function with given fields: exchange, routingKeys, queue, queueOptions, queueTimeout
func (_m *AMQPConsumer) Receive(exchange string, routingKeys []string, queue string, queueOptions simpleamqp.QueueOptions, queueTimeout time.Duration) chan simpleamqp.AmqpMessage {
	ret := _m.Called(exchange, routingKeys, queue, queueOptions, queueTimeout)

	var r0 chan simpleamqp.AmqpMessage
	if rf, ok := ret.Get(0).(func(string, []string, string, simpleamqp.QueueOptions, time.Duration) chan simpleamqp.AmqpMessage); ok {
		r0 = rf(exchange, routingKeys, queue, queueOptions, queueTimeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan simpleamqp.AmqpMessage)
		}
	}

	return r0
}

// ReceiveWithoutTimeout provides a mock function with given fields: exchange, routingKeys, queue, queueOptions
func (_m *AMQPConsumer) ReceiveWithoutTimeout(exchange string, routingKeys []string, queue string, queueOptions simpleamqp.QueueOptions) chan simpleamqp.AmqpMessage {
	ret := _m.Called(exchange, routingKeys, queue, queueOptions)

	var r0 chan simpleamqp.AmqpMessage
	if rf, ok := ret.Get(0).(func(string, []string, string, simpleamqp.QueueOptions) chan simpleamqp.AmqpMessage); ok {
		r0 = rf(exchange, routingKeys, queue, queueOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan simpleamqp.AmqpMessage)
		}
	}

	return r0
}
