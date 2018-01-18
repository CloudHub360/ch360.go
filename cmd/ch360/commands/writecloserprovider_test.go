package commands_test

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/CloudHub360/ch360.go/io_util"
	iomocks "github.com/CloudHub360/ch360.go/io_util/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WriteCloserProviderSuite struct {
	suite.Suite
	mockWriteCloser         *iomocks.WriteCloser
	mockWriterProvider      *mocks.WriterProvider
	mockWriteCloserProvider *mocks.WriteCloserProvider
}

func (suite *WriteCloserProviderSuite) TearDownTest() {}

func (suite *WriteCloserProviderSuite) SetupTest() {
	suite.mockWriteCloser = &iomocks.WriteCloser{}
	suite.mockWriterProvider = &mocks.WriterProvider{}
	suite.mockWriteCloserProvider = &mocks.WriteCloserProvider{}
}

func TestWriteCloserProviderSuiteRunner(t *testing.T) {
	suite.Run(t, new(WriteCloserProviderSuite))
}

func (suite *WriteCloserProviderSuite) Test_BasicWriterFactory_Returns_Its_Param() {
	// Arrange
	expectedWriter := &bytes.Buffer{}
	sut := commands.NewBasicWriterFactory(expectedWriter)

	// Act
	receivedWriterCloser, err := sut.Provide(generators.String("fullpath"))

	// Assert
	assert.Equal(suite.T(), expectedWriter, receivedWriterCloser)
	assert.Nil(suite.T(), err)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Nil_For_Empty_Param() {
	// Arrange
	sut := commands.NewAutoClosingWriterFactory(suite.mockWriteCloserProvider)

	// Act
	receivedWriter, _ := sut.Provide("")

	// Assert
	assert.Nil(suite.T(), receivedWriter)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Calls_Underlying() {
	// Arrange
	sut := commands.NewAutoClosingWriterFactory(suite.mockWriteCloserProvider)
	suite.mockWriteCloserProvider.On("Provide", mock.Anything).Return(nil, nil)

	// Act
	path := generators.String("fullpath")
	sut.Provide(path)

	// Assert
	suite.mockWriteCloserProvider.AssertCalled(suite.T(), "Provide", path)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Error_From_Underlying() {
	// Arrange
	sut := commands.NewAutoClosingWriterFactory(suite.mockWriteCloserProvider)
	expectedError := errors.New("simulated error")
	suite.mockWriteCloserProvider.On("Provide", mock.Anything).Return(nil, expectedError)

	// Act
	_, receivedError := sut.Provide(generators.String("fullpath"))

	// Assert
	assert.Equal(suite.T(), expectedError, receivedError)
}

func (suite *WriteCloserProviderSuite) Test_AutoClosingWriteCloserProvider_Returns_Provided_Underlying_Within_AutoCloser() {
	// Arrange
	sut := commands.NewAutoClosingWriterFactory(suite.mockWriteCloserProvider)

	suite.mockWriteCloserProvider.On("Provide", mock.Anything).Return(suite.mockWriteCloser, nil)

	// Act
	receivedWriter, _ := sut.Provide(generators.String("fullpath"))

	// Assert
	ac, ok := receivedWriter.(*io_util.AutoCloser)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), suite.mockWriteCloser, ac.Underlying)
}
