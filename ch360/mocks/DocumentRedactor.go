// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import io "io"
import mock "github.com/stretchr/testify/mock"
import request "github.com/waives/surf/ch360/request"

// DocumentRedactor is an autogenerated mock type for the DocumentRedactor type
type DocumentRedactor struct {
	mock.Mock
}

// Redact provides a mock function with given fields: ctx, documentId, redactRequest
func (_m *DocumentRedactor) Redact(ctx context.Context, documentId string, redactRequest request.RedactedPdfRequest) (io.ReadCloser, error) {
	ret := _m.Called(ctx, documentId, redactRequest)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context, string, request.RedactedPdfRequest) io.ReadCloser); ok {
		r0 = rf(ctx, documentId, redactRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, request.RedactedPdfRequest) error); ok {
		r1 = rf(ctx, documentId, redactRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
