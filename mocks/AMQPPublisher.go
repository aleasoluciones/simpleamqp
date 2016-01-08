package mocks

import "github.com/stretchr/testify/mock"

type AMQPPublisher struct {
	mock.Mock
}

func (_m *AMQPPublisher) Publish(_a0 string, _a1 []byte) {
	_m.Called(_a0, _a1)
}
func (_m *AMQPPublisher) PublishWithTTL(_a0 string, _a1 []byte, _a2 int) {
	_m.Called(_a0, _a1, _a2)
}
