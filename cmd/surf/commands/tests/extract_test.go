package tests

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type extractCommandSuite struct {
	suite.Suite
	extractionService *mocks.ExtractionService
	sut               *commands.ExtractCmd
	filePatterns      []string
	ctx               context.Context
	expectedErr       error
	extractorName     string
}

func (suite *extractCommandSuite) SetupTest() {
	suite.extractionService = new(mocks.ExtractionService)
	suite.ctx, _ = context.WithCancel(context.Background())
	suite.filePatterns = []string{generators.String("file1"), generators.String("file2")}
	suite.extractorName = generators.String("extractorName")
	suite.expectedErr = errors.New("simulated error")

	suite.extractionService.
		On("ExtractAll", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.expectedErr)

	suite.sut = &commands.ExtractCmd{
		FilePaths:         suite.filePatterns,
		ExtractionService: suite.extractionService,
		ExtractorName:     suite.extractorName,
	}
}

func TestExtractCommandSuiteRunner(t *testing.T) {
	suite.Run(t, new(extractCommandSuite))
}

func (suite *extractCommandSuite) Test_ExtractionService_ExtractAll_Called_With_Correct_Params() {
	_ = suite.sut.Execute(suite.ctx)

	suite.extractionService.AssertCalled(suite.T(),
		"ExtractAll",
		suite.ctx,
		suite.filePatterns,
		suite.extractorName)
}

func (suite *extractCommandSuite) Test_Error_Returned_From_ExtractionService() {
	actualErr := suite.sut.Execute(suite.ctx)

	assert.EqualError(suite.T(), actualErr, suite.expectedErr.Error())
}
