package tests

import (
	"github.com/CloudHub360/ch360.go/ch360/types"
	formatterMocks "github.com/CloudHub360/ch360.go/output/formatters/mocks"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	sinkMocks "github.com/CloudHub360/ch360.go/output/sinks/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CombinedResultsWriterSuite struct {
	suite.Suite
	sut                  *resultsWriters.CombinedResultsWriter
	sink                 *sinkMocks.Sink
	formatter            *formatterMocks.ClassifyResultsFormatter
	filename             string
	classificationResult *types.ClassificationResult
}

func (suite *CombinedResultsWriterSuite) SetupTest() {
	suite.sink = new(sinkMocks.Sink)
	suite.sink.On("Open").Return(nil)
	suite.sink.On("Close").Return(nil)

	suite.formatter = new(formatterMocks.ClassifyResultsFormatter)
	suite.formatter.On("WriteHeader", mock.Anything).Return(nil)
	suite.formatter.On("WriteSeparator", mock.Anything).Return(nil)
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.formatter.On("WriteFooter", mock.Anything).Return(nil)

	suite.sut = resultsWriters.NewCombinedResultsWriter(suite.sink, suite.formatter)

	suite.filename = generators.String("filename")
	suite.classificationResult = &types.ClassificationResult{
		DocumentType: generators.String("documentType"),
		IsConfident:  generators.Bool(),
	}
}

func TestCombinedResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CombinedResultsWriterSuite))
}

func (suite *CombinedResultsWriterSuite) TestStart_OpensSink_Then_Writes_Header() {
	err := suite.sut.Start()

	assert.Nil(suite.T(), err)
	//TODO Fake sink to check opened first?
	suite.sink.AssertCalled(suite.T(), "Open")
	suite.formatter.AssertCalled(suite.T(), "WriteHeader", suite.sink)
}

func (suite *CombinedResultsWriterSuite) TestStart_Returns_Error_From_WriteHeader() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteHeader", mock.Anything).Return(expectedErr)

	err := suite.sut.Start()
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestStart_Returns_Error_From_SinkOpen() {
	suite.sink.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.sink.On("Open").Return(expectedErr)

	err := suite.sut.Start()
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Writes_Result_To_Sink_But_No_Separator_On_First_Call() {
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Nil(suite.T(), err)
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, suite.filename, suite.classificationResult)
	suite.formatter.AssertNotCalled(suite.T(), "WriteSeparator", suite.sink)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Writes_Separator_Then_Result_On_Subsequent_Calls() {
	filename1 := generators.String("filename1")
	filename2 := generators.String("filename2")
	err := suite.sut.WriteResult(filename1, suite.classificationResult)
	err = suite.sut.WriteResult(filename2, suite.classificationResult)

	assert.Nil(suite.T(), err)
	require.Len(suite.T(), suite.formatter.Calls, 3)
	// Check calls were in the correct order
	assert.Equal(suite.T(), "WriteResult", suite.formatter.Calls[0].Method)
	assert.Equal(suite.T(), "WriteSeparator", suite.formatter.Calls[1].Method)
	assert.Equal(suite.T(), "WriteResult", suite.formatter.Calls[2].Method)
	// Check calls had the correct parameters
	suite.formatter.AssertCalled(suite.T(), "WriteSeparator", suite.sink)
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, filename1, suite.classificationResult)
	suite.formatter.AssertCalled(suite.T(), "WriteResult", suite.sink, filename2, suite.classificationResult)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Returns_Error_From_WriteResult() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestWriteResult_Returns_Error_From_WriteSeparator() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteSeparator", mock.Anything).Return(expectedErr)
	suite.formatter.On("WriteResult", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	suite.sut.WriteResult(suite.filename, suite.classificationResult)
	err := suite.sut.WriteResult(suite.filename, suite.classificationResult)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CombinedResultsWriterSuite) TestFinish_Writes_Footer_Then_ClosesSink() {
	err := suite.sut.Finish()

	assert.Nil(suite.T(), err)
	//TODO Fake sink to check WriteFooter first?
	suite.sink.AssertCalled(suite.T(), "Close")
	suite.formatter.AssertCalled(suite.T(), "WriteFooter", suite.sink)
}

func (suite *CombinedResultsWriterSuite) TestFinish_Returns_Error_From_WriteFooter() {
	suite.formatter.ExpectedCalls = nil
	expectedErr := errors.New("expectedError")
	suite.formatter.On("WriteFooter", mock.Anything).Return(expectedErr)

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
