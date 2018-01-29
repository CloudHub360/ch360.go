package tests

import (
	"bytes"
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateExtractorSuite struct {
	suite.Suite
	output  *bytes.Buffer
	creator *mocks.ExtractorCreator
	sut     *commands.CreateExtractor
	config  *bytes.Buffer
}

func (suite *CreateExtractorSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.creator = new(mocks.ExtractorCreator)

	suite.sut = commands.NewCreateExtractor(suite.output, suite.creator)
	suite.config = &bytes.Buffer{}
	suite.config.Write([]byte("some data"))

	suite.creator.On("Create", mock.Anything, mock.Anything).Return(nil)
}

func TestCreateExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateExtractorSuite))
}

func (suite *CreateExtractorSuite) ClearExpectedCalls() {
	suite.creator.ExpectedCalls = nil
}

func (suite *CreateExtractorSuite) TestCreateExtractor_Execute_Calls_Client_With_Correct_Args() {
	suite.sut.Execute("charlie", suite.config)

	suite.creator.AssertCalled(suite.T(), "Create", "charlie", suite.config)
}

func (suite *CreateExtractorSuite) TestCreateExtractor_Execute_Returns_Error_If_The_Extractor_Cannot_Be_Created() {
	expectedErr := errors.New("Error message")
	suite.ClearExpectedCalls()
	suite.creator.On("Create", mock.Anything, mock.Anything).Return(expectedErr)

	receivedErr := suite.sut.Execute("charlie", suite.config)

	assert.Equal(suite.T(), expectedErr, receivedErr)
}
