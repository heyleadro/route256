// Code generated by mockery v2.30.1. DO NOT EDIT.

package checkout

import (
	context "context"
	model "route256/checkout/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// MockCartLister is an autogenerated mock type for the CartLister type
type MockCartLister struct {
	mock.Mock
}

// GetProducts provides a mock function with given fields: ctx, sku
func (_m *MockCartLister) GetProducts(ctx context.Context, sku uint32) (model.Item, error) {
	ret := _m.Called(ctx, sku)

	var r0 model.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (model.Item, error)); ok {
		return rf(ctx, sku)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) model.Item); ok {
		r0 = rf(ctx, sku)
	} else {
		r0 = ret.Get(0).(model.Item)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, sku)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockCartLister creates a new instance of MockCartLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCartLister(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCartLister {
	mock := &MockCartLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
