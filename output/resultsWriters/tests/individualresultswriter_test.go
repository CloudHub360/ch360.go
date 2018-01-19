package tests

import (
	"errors"
	"github.com/CloudHub360/ch360.go/ch360/types"
	formatterMocks "github.com/CloudHub360/ch360.go/output/formatters/mocks"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	sinkMocks "github.com/CloudHub360/ch360.go/output/sinks/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IndividualResultsWriterSuite struct {
	suite.Suite
	sut                  *resultsWriters.IndividualResultsWriter
	sink                 *sinkMocks.Sink
	sinkFactory          *sinkMocks.SinkFactory
	formatter            *formatterMocks.ClassifyResultsFormatter
	filename             string
	classificationResult *types.ClassificationResult
}

func (suite *IndividualResultsWriterSuite) SetupTest() {
	suite.sink = new(sinkMocks.Sink)
	suite.sink.On("Open").Return(nil)
	suite.sink.On("Close").Return(nil)

	suite.sinkFactory = new(sinkMocks.SinkFactory)
	suite.sinkFactory.On("Sink", mock.Anything).Return(suite.sink, nil)

	suite.formatter = new(formatterMocks.ClassifyResultsFormatter)
	suite.formatter.On("WriteHeader", mock.Anything).Return(nil)
	suite.formatter.On("WriteSeparator", mock.Anything).Return(nil)
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.formatter.On("WriteFooter", mock.Anything).Return(nil)

	suite.sut = resultsWriters.NewIndividualResultsWriter(suite.sinkFactory, suite.formatter)

	suite.filename = generators.String("filename")
	suite.classificationResult = &types.ClassificationResult{
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
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, suite.filename, suite.classificationResult)
}

func (suite *IndividualResultsWriterSuite) TestWriteResults_Returns_Error_From_WriteResult() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

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
