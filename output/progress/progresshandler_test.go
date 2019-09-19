package progress_test

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/output/progress"
	"github.com/waives/surf/output/resultsWriters/mocks"
	"github.com/waives/surf/test/generators"
	"math/rand"
	"testing"
	"time"
)

type ClassifyProgressHandlerSuite struct {
	suite.Suite
	suts             []*progress.ProgressHandler
	mockResultWriter *mocks.ResultsWriter
	outBuffer        *bytes.Buffer
}

func (suite *ClassifyProgressHandlerSuite) TearDownTest() {
	for _, sut := range suite.suts {
		_ = sut.NotifyFinish()
	}
}

func (suite *ClassifyProgressHandlerSuite) SetupTest() {
	suite.mockResultWriter = &mocks.ResultsWriter{}
	suite.outBuffer = new(bytes.Buffer)
	suite.setupMockResultWriter(true, true, true)

	rand.Seed(time.Now().Unix())

	suite.suts = []*progress.ProgressHandler{
		progress.NewProgressHandler(suite.mockResultWriter, true, suite.outBuffer),
		progress.NewProgressHandler(suite.mockResultWriter, false, suite.outBuffer),
	}
}

func (suite *ClassifyProgressHandlerSuite) setupMockResultWriter(start, write, finish bool) {
	suite.mockResultWriter.ExpectedCalls = nil

	if start {
		suite.mockResultWriter.On("Start").Return(nil)
	}

	if write {
		suite.mockResultWriter.On("WriteResult", mock.Anything, mock.Anything).Return(nil)
	}

	if finish {
		suite.mockResultWriter.On("Finish").Return(nil)
	}
}

func TestClassifyProgressHandlerSuiteRunner(t *testing.T) {
	suite.Run(t, new(ClassifyProgressHandlerSuite))
}

func (suite *ClassifyProgressHandlerSuite) Test_ClassifyProgressHandler_Calls_Underlying_ResultWriter() {
	for _, sut := range suite.suts {
		// Arrange
		expectedFilename := generators.String("filename")
		expectedResult := AClassificationResult()

		// Act
		sut.NotifyStart(1)
		sut.Notify(expectedFilename, expectedResult)

		// Assert
		suite.mockResultWriter.AssertCalled(suite.T(), "WriteResult", expectedFilename, expectedResult)
	}
}

func (suite *ClassifyProgressHandlerSuite) Test_ClassifyProgressHandler_Returns_Error_From_ResultWriter() {
	for _, sut := range suite.suts {
		// Arrange
		expectedErr := errors.New("simulated error")
		suite.setupMockResultWriter(true, false, true)
		suite.mockResultWriter.On("WriteResult", mock.Anything, mock.Anything).Return(expectedErr)

		// Act
		sut.NotifyStart(rand.Int())
		receivedErr := sut.Notify(generators.String("filename"), AClassificationResult())

		// Assert
		suite.Assert().Equal(expectedErr, receivedErr)
	}
}

func (suite *ClassifyProgressHandlerSuite) Test_ClassifyProgressHandler_Returns_Error_If_Notify_Is_Called_Before_NotifyStart() {
	for _, sut := range suite.suts {
		// Act
		err := sut.Notify(generators.String("filename"), AClassificationResult())

		// Assert
		suite.Assert().NotNil(err)
		suite.mockResultWriter.AssertNotCalled(suite.T(), "WriteResult", mock.Anything, mock.Anything)
	}

}

func (suite *ClassifyProgressHandlerSuite) Test_ClassifyProgressHandler_Returns_Error_If_NotifyErr_Is_Called_Before_NotifyStart() {
	for _, sut := range suite.suts {
		// Act
		err := sut.NotifyErr(generators.String("filename"), errors.New("simulated error"))

		// Assert
		suite.Assert().NotNil(err)
		suite.mockResultWriter.AssertNotCalled(suite.T(), "WriteResult", mock.Anything, mock.Anything)
	}
}

func (suite *ClassifyProgressHandlerSuite) Test_ClassifyProgressHandler_Returns_Error_If_NotifyFinish_Is_Called_Before_NotifyStart() {
	for _, sut := range suite.suts {
		// Act
		err := sut.NotifyFinish()

		// Assert
		suite.Assert().NotNil(err)
		suite.mockResultWriter.AssertNotCalled(suite.T(), "WriteResult", mock.Anything, mock.Anything)
	}
}

func AClassificationResult() *results.ClassificationResult {
	return &results.ClassificationResult{}
}
