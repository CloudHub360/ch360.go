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

type classifyCommandSuite struct {
	suite.Suite
	classificationService *mocks.ClassificationService
	sut                   *commands.ClassifyCmd
	filePatterns          []string
	ctx                   context.Context
	expectedErr           error
	classifierName        string
}

func (suite *classifyCommandSuite) SetupTest() {
	suite.classificationService = new(mocks.ClassificationService)
	suite.ctx, _ = context.WithCancel(context.Background())
	suite.filePatterns = []string{generators.String("file1"), generators.String("file2")}
	suite.classifierName = generators.String("classifierName")
	suite.expectedErr = errors.New("simulated error")

	suite.classificationService.
		On("ClassifyAll", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.expectedErr)

	suite.sut = &commands.ClassifyCmd{
		FilePaths:             suite.filePatterns,
		ClassificationService: suite.classificationService,
		ClassifierName:        suite.classifierName,
	}
}

func TestClassifyCommandSuiteRunner(t *testing.T) {
	suite.Run(t, new(classifyCommandSuite))
}

func (suite *classifyCommandSuite) Test_ClassificationService_ClassifyAll_Called_With_Correct_Params() {
	_ = suite.sut.Execute(suite.ctx)

	suite.classificationService.AssertCalled(suite.T(),
		"ClassifyAll",
		suite.ctx,
		suite.filePatterns,
		suite.classifierName)
}

func (suite *classifyCommandSuite) Test_Error_Returned_From_ClassificationService() {
	actualErr := suite.sut.Execute(suite.ctx)

	assert.EqualError(suite.T(), actualErr, suite.expectedErr.Error())
}
