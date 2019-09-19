package tests

import (
	"errors"
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

type IndividualResultsWriterSuite struct {
	suite.Suite
	sut                  *resultsWriters.IndividualResultsWriter
	sink                 *sinkMocks.Sink
	sinkFactory          *sinkMocks.SinkFactory
	formatter            *formatterMocks.ResultsFormatter
	filename             string
	classificationResult *results.ClassificationResult
}

func (suite *IndividualResultsWriterSuite) SetupTest() {
	suite.sink = new(sinkMocks.Sink)
	suite.sink.On("Open").Return(nil)
	suite.sink.On("Close").Return(nil)

	suite.sinkFactory = new(sinkMocks.SinkFactory)
	suite.sinkFactory.On("Sink", mock.Anything).Return(suite.sink, nil)

	suite.formatter = new(formatterMocks.ResultsFormatter)
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.formatter.On("Flush", mock.Anything).Return(nil)

	suite.sut = resultsWriters.NewIndividualResultsWriter(suite.sinkFactory, suite.formatter)

	suite.filename = generators.String("filename")
	suite.classificationResult = &results.ClassificationResult{
		DocumentType: generators.String("documentType"),
		IsConfident:  generators.Bool(),
	}
}

func TestIndividualResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(IndividualResultsWriterSuite))
}

func (suite *IndividualResultsWriterSuite) TestStart_Does_Nothing() {
	err := suite.sut.Start()

	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.sinkFactory.Calls, 0)
	assert.Len(suite.T(), suite.formatter.Calls, 0)
}

func (suite *IndividualResultsWriterSuite) TestFinish_Does_Nothing() {
	err := suite.sut.Finish()

	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.sinkFactory.Calls, 0)
	assert.Len(suite.T(), suite.formatter.Calls, 0)
}

func (suite *IndividualResultsWriterSuite) TestWriteResult_Gets_Sink_From_SinkFactory() {
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Nil(suite.T(), err)
	suite.sinkFactory.AssertCalled(suite.T(), "Sink", mock.Anything)
}

func (suite *IndividualResultsWriterSuite) TestWriteResult_Returns_Error_From_SinkFactory() {
	suite.sinkFactory.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.sinkFactory.On("Sink", mock.Anything).Return(nil, expectedErr)

	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_OpensSink() {
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Nil(suite.T(), err)
	suite.sink.AssertCalled(suite.T(), "Open")
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_Writes_Results() {
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Nil(suite.T(), err)
	suite.formatter.AssertCalled(suite.T(), "WriteResult",
		suite.sink,
		suite.filename,
		suite.classificationResult,
		formatters.IncludeHeader)
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_Returns_Error_From_WriteResult() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_ClosesSink() {
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Nil(suite.T(), err)
	suite.sink.AssertCalled(suite.T(), "Close")
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_Returns_Error_From_CloseSink() {
	suite.sink.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.sink.On("Open").Return(nil)
	suite.sink.On("Close").Return(expectedErr)

	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}
