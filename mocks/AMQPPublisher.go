// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"

// AMQPPublisher is an autogenerated mock type for the AMQPPublisher type
type AMQPPublisher struct {
	mock.Mock
}

// Publish provides a mock function with given fields: _a0, _a1
func (_m *AMQPPublisher) Publish(_a0 string, _a1 []byte) {
	_m.Called(_a0, _a1)
}

// PublishWithTTL provides a mock function with given fields: _a0, _a1, _a2
func (_m *AMQPPublisher) PublishWithTTL(_a0 string, _a1 []byte, _a2 int) {
	_m.Called(_a0, _a1, _a2)
}
