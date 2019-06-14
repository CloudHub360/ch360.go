package tests

import (
	"bytes"
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

type UploadExtractorSuite struct {
	suite.Suite
	output        *bytes.Buffer
	creator       *mocks.ExtractorCreator
	sut           *commands.UploadExtractor
	config        *bytes.Buffer
	extractorName string
	ctx           context.Context
}

func (suite *UploadExtractorSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.creator = new(mocks.ExtractorCreator)

	suite.config = &bytes.Buffer{}
	suite.config.Write([]byte("some data"))
	suite.extractorName = generators.String("extractor-name")
	suite.sut = commands.NewUploadExtractor(suite.output,
		suite.creator,
		suite.extractorName,
		suite.config)

	suite.creator.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.ctx = context.Background()
}

func TestUploadExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(UploadExtractorSuite))
}

func (suite *UploadExtractorSuite) ClearExpectedCalls() {
	suite.creator.ExpectedCalls = nil
}

func (suite *UploadExtractorSuite) TestUploadExtractor_Execute_Calls_Client_With_Correct_Args() {
	suite.sut.Execute(context.Background())

	suite.creator.AssertCalled(suite.T(), "Create", suite.ctx, suite.extractorName, suite.config)
}

func (suite *UploadExtractorSuite) TestUploadExtractor_Execute_Returns_Error_If_The_Extractor_Cannot_Be_Created() {
	expectedErr := errors.New("Error message")
	suite.ClearExpectedCalls()
	suite.creator.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	receivedErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}
