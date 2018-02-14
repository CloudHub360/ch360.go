// Code generated by mockery v1.0.0
package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import results "github.com/CloudHub360/ch360.go/ch360/results"

// DocumentExtractor is an autogenerated mock type for the DocumentExtractor type
type DocumentExtractor struct {
	mock.Mock
}

// Extract provides a mock function with given fields: ctx, documentId, extractorName
func (_m *DocumentExtractor) Extract(ctx context.Context, documentId string, extractorName string) (*results.ExtractionResult, error) {
	ret := _m.Called(ctx, documentId, extractorName)

	var r0 *results.ExtractionResult
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *results.ExtractionResult); ok {
		r0 = rf(ctx, documentId, extractorName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*results.ExtractionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, documentId, extractorName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
