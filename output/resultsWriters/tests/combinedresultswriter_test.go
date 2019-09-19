package tests

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/output/formatters"
	formatterMocks "github.com/waives/surf/output/formatters/mocks"
	"github.com/waives/surf/output/resultsWriters"
	sinkMocks "github.com/waives/surf/output/sinks/mocks"
	"github.com/waives/surf/test/generators"
	"testing"
)

type CombinedResultsWriterSuite struct {
	suite.Suite
	sut                  *resultsWriters.CombinedResultsWriter
	sink                 *sinkMocks.Sink
	formatter            *formatterMocks.ResultsFormatter
	filename             string
	classificationResult *results.ClassificationResult
}

func (suite *CombinedResultsWriterSuite) SetupTest() {
	suite.sink = new(sinkMocks.Sink)
	suite.sink.On("Open").Return(nil)
	suite.sink.On("Close").Return(nil)

	suite.formatter = new(formatterMocks.ResultsFormatter)
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.formatter.On("Flush", mock.Anything).Return(nil)

	suite.sut = resultsWriters.NewCombinedResultsWriter(suite.sink, suite.formatter)

	suite.filename = generators.String("filename")
	suite.classificationResult = &results.ClassificationResult{
		DocumentType: generators.String("documentType"),
		IsConfident:  generators.Bool(),
	}
}

func TestCombinedResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CombinedResultsWriterSuite))
}

func (suite *CombinedResultsWriterSuite) TestStart_OpensSink() {
	err := suite.sut.Start()

	assert.Nil(suite.T(), err)
	suite.sink.AssertCalled(suite.T(), "Open")
}

func (suite *CombinedResultsWriterSuite) TestStart_OpensSink_Before_Writing() {
	sink := &fakeSink{}
	sut := resultsWriters.NewCombinedResultsWriter(sink, formatters.NewCSVClassifyResultsFormatter())
	err := sut.Start()

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), sink.IsOpen)
}

func (suite *CombinedResultsWriterSuite) TestStart_Returns_Error_From_SinkOpen() {
	suite.sink.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.sink.On("Open").Return(expectedErr)

	err := suite.sut.Start()
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Returns_Error_From_WriteResult() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestFinish_Calls_Flush_Then_ClosesSink() {
	err := suite.sut.Finish()

	assert.Nil(suite.T(), err)
	suite.sink.AssertCalled(suite.T(), "Close")
	suite.formatter.AssertCalled(suite.T(), "Flush", suite.sink)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Writes_Header_On_First_Call_Only() {
	suite.sut.WriteResult(suite.filename, suite.classificationResult)
	suite.sut.WriteResult(suite.filename, suite.classificationResult)

	suite.formatter.AssertNumberOfCalls(suite.T(), "WriteResult", 2)
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, suite.filename, suite.classificationResult, formatters.IncludeHeader)
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, suite.filename, suite.classificationResult, formatters.FormatOption(0))
}

func (suite *CombinedResultsWriterSuite) TestFinish_Returns_Error_From_Flush() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("Flush", mock.Anything).Return(expectedErr)

	err := suite.sut.Finish()
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestFinish_Returns_Error_From_SinkClose() {
	suite.sink.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.sink.On("Close").Return(expectedErr)

	err := suite.sut.Finish()
	assert.Equal(suite.T(), expectedErr, err)
}

type fakeSink struct {
	IsOpen bool
}

func (f *fakeSink) Open() error {
	f.IsOpen = true
	return nil
}

func (f *fakeSink) Close() error {
	return nil
}

func (f *fakeSink) Write(b []byte) (int, error) {
	if f.IsOpen {
		return 0, nil
	} else {
		return 0, errors.New("Sink hasn't been opened")
	}
}
