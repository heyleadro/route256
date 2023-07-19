// Code generated by mockery v2.30.1. DO NOT EDIT.

package loms

import (
	model "route256/loms/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// MockSender is an autogenerated mock type for the Sender type
type MockSender struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: message
func (_m *MockSender) SendMessage(message model.ProducerMessage) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.ProducerMessage) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockSender creates a new instance of MockSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSender {
	mock := &MockSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}