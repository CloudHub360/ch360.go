// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"

// ClassifierDeleter is an autogenerated mock type for the ClassifierDeleter type
type ClassifierDeleter struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, name
func (_m *ClassifierDeleter) Delete(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
