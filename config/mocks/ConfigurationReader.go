// Code generated by mockery v1.0.0
package mocks

import config "github.com/CloudHub360/ch360.go/config"
import mock "github.com/stretchr/testify/mock"

// ConfigurationReader is an autogenerated mock type for the ConfigurationReader type
type ConfigurationReader struct {
	mock.Mock
}

// ReadConfiguration provides a mock function with given fields:
func (_m *ConfigurationReader) ReadConfiguration() (*config.Configuration, error) {
	ret := _m.Called()

	var r0 *config.Configuration
	if rf, ok := ret.Get(0).(func() *config.Configuration); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Configuration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}