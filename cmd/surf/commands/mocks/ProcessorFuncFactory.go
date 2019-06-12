// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import pool "github.com/CloudHub360/ch360.go/pool"

// ProcessorFuncFactory is an autogenerated mock type for the ProcessorFuncFactory type
type ProcessorFuncFactory struct {
	mock.Mock
}

// ProcessorFor provides a mock function with given fields: ctx, filename
func (_m *ProcessorFuncFactory) ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc {
	ret := _m.Called(ctx, filename)

	var r0 pool.ProcessorFunc
	if rf, ok := ret.Get(0).(func(context.Context, string) pool.ProcessorFunc); ok {
		r0 = rf(ctx, filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pool.ProcessorFunc)
		}
	}

	return r0
}
