package io_util_test

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"github.com/CloudHub360/ch360.go/io_util/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type AutoCloserSuite struct {
	suite.Suite
	sut        *io_util.AutoCloser
	underlying *mocks.WriteCloser
}

func (suite *AutoCloserSuite) SetupTest() {
	suite.underlying = &mocks.WriteCloser{}
	suite.sut = &io_util.AutoCloser{
		Underlying: suite.underlying,
	}

	rand.Seed(time.Now().Unix())
}

func TestAutoCloserSuiteRunner(t *testing.T) {
	suite.Run(t, new(AutoCloserSuite))
}

func (suite *AutoCloserSuite) Test_AutoCloser_Calls_Close_After_Write() {
	// Arrange
	suite.underlying.On("Close").Return(nil)
	suite.underlying.On("Write", mock.Anything).Return(0, nil)

	// Act
	suite.sut.Write(nil)

	// Assert
	suite.underlying.AssertCalled(suite.T(), "Close")
}

func (suite *AutoCloserSuite) Test_AutoCloser_Returns_Values_From_Write() {
	// Arrange
	var (
		expectedByteCount = rand.Int()
		expectedErr       = errors.New("simulated error")
	)

	suite.underlying.On("Close").Return(nil)
	suite.underlying.On("Write", mock.Anything).Return(expectedByteCount, expectedErr)

	// Act
	receivedByteCount, receivedErr := suite.sut.Write(nil)

	// Assert
	assert.Equal(suite.T(), expectedByteCount, receivedByteCount)
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *AutoCloserSuite) Test_AutoCloser_Returns_Error_From_Close() {
	// Arrange
	var (
		expectedErr = errors.New("simulated error")
	)

	suite.underlying.On("Close").Return(expectedErr)
	suite.underlying.On("Write", mock.Anything).Return(0, nil)

	// Act
	_, receivedErr := suite.sut.Write(nil)

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}
