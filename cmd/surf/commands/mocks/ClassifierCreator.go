// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"

// ClassifierCreator is an autogenerated mock type for the ClassifierCreator type
type ClassifierCreator struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, name
func (_m *ClassifierCreator) Create(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
