package io_util_test

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AutoCloserSuite struct {
	suite.Suite
	sut *io_util.AutoCloser
}

func (suite *AutoCloserSuite) SetupTest() {

}

func TestAutoCloserSuiteRunner(t *testing.T) {
	suite.Run(t, new(AutoCloserSuite))
}

func (suite *AutoCloserSuite) Test_AutoCloser_Calls_Close_After_Write() {

}
