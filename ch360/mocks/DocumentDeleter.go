// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"

// DocumentDeleter is an autogenerated mock type for the DocumentDeleter type
type DocumentDeleter struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, documentId
func (_m *DocumentDeleter) Delete(ctx context.Context, documentId string) error {
	ret := _m.Called(ctx, documentId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, documentId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
