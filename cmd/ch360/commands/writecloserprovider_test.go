package commands_test

import (
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/io_util"
	"github.com/CloudHub360/ch360.go/io_util/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

type WriteCloserProviderSuite struct {
	suite.Suite
	mockWriteCloser     *mocks.WriteCloser
	writeCloserProvider commands.WriteCloserProvider
}

func (suite *WriteCloserProviderSuite) TearDownTest() {}

func (suite *WriteCloserProviderSuite) SetupTest() {
	suite.mockWriteCloser = &mocks.WriteCloser{}
	suite.writeCloserProvider = func(fullPath string) (io.WriteCloser, error) {
		return nil, nil
	}
}

func TestWriteCloserProviderSuiteRunner(t *testing.T) {
	suite.Run(t, new(WriteCloserProviderSuite))
}

func (suite *WriteCloserProviderSuite) Test_DummyProvider_Returns_Its_Param() {
	// Arrange
	sut := commands.NewDummyWriteCloserProvider

	// Act
	receivedWriterCloser, err := sut(suite.mockWriteCloser)(generators.String("fullpath"))

	// Assert
	assert.Equal(suite.T(), suite.mockWriteCloser, receivedWriterCloser)
	assert.Nil(suite.T(), err)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Nil_For_Empty_Param() {
	// Arrange
	sut := commands.NewAutoClosingWriteCloserProvider

	// Act
	receivedWriterCloser, _ := sut(suite.writeCloserProvider)("")

	// Assert
	assert.Nil(suite.T(), receivedWriterCloser)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Calls_Underlying() {
	// Arrange
	sut := commands.NewAutoClosingWriteCloserProvider
	underlyingCalled := false
	var writeCloserProvider commands.WriteCloserProvider = func(fullPath string) (io.WriteCloser, error) {
		underlyingCalled = true
		return nil, nil
	}

	// Act
	sut(writeCloserProvider)(generators.String("fullpath"))

	// Assert
	assert.True(suite.T(), underlyingCalled)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Error_From_Underlying() {
	// Arrange
	sut := commands.NewAutoClosingWriteCloserProvider
	expectedError := errors.New("simulated error")
	var writeCloserProvider commands.WriteCloserProvider = func(fullPath string) (io.WriteCloser, error) {
		return nil, expectedError
	}

	// Act
	_, receivedError := sut(writeCloserProvider)(generators.String("fullpath"))

	// Assert
	assert.Equal(suite.T(), expectedError, receivedError)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Provided_Underlying_Within_AutoCloser() {
	// Arrange
	sut := commands.NewAutoClosingWriteCloserProvider
	expectedWriteCloser := &mocks.WriteCloser{}
	var writeCloserProvider commands.WriteCloserProvider = func(fullPath string) (io.WriteCloser, error) {
		return expectedWriteCloser, nil
	}

	// Act
	receivedProvider, _ := sut(writeCloserProvider)(generators.String("fullpath"))

	// Assert
	ac, ok := receivedProvider.(*io_util.AutoCloser)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), expectedWriteCloser, ac.Underlying)
}
