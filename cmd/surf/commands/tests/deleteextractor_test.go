package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/cmd/surf/commands"
	"github.com/waives/surf/cmd/surf/commands/mocks"
	"testing"
)

type DeleteExtractorSuite struct {
	suite.Suite
	sut           *commands.DeleteExtractorCmd
	output        *bytes.Buffer
	client        *mocks.ExtractorDeleterGetter
	extractorName string
	ctx           context.Context
}

func (suite *DeleteExtractorSuite) SetupTest() {
	suite.client = new(mocks.ExtractorDeleterGetter)
	suite.client.On("GetAll", mock.Anything).Return(
		AListOfExtractors("charlie", "jo", "chris"), nil)

	suite.client.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.output = &bytes.Buffer{}

	suite.extractorName = "charlie"

	suite.sut = &commands.DeleteExtractorCmd{
		Client:        suite.client,
		ExtractorName: suite.extractorName,
	}
	suite.ctx = context.Background()
}

func TestDeleteExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(DeleteExtractorSuite))
}

func (suite *DeleteExtractorSuite) TestDeleteExtractor_Execute_Deletes_The_Named_Extractor_When_It_Exists() {
	suite.sut.Execute(context.Background())

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
	suite.client.AssertCalled(suite.T(), "Delete", suite.ctx, "charlie")
}

func (suite *DeleteExtractorSuite) TestDeleteExtractor_Execute_Does_Not_Delete_The_Named_Extractor_When_It_Does_Not_Exist() {
	suite.sut.Execute(context.Background())

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
	suite.client.AssertNotCalled(suite.T(), "Delete")
}

func (suite *DeleteExtractorSuite) TestDeleteExtractor_Execute_Returns_An_Error_If_The_Extractors_Cannot_Be_Retrieved() {
	suite.client.ExpectedCalls = nil
	expected := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expected)

	actual := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expected, actual)
	suite.client.AssertNotCalled(suite.T(), "Delete")
}

func (suite *DeleteExtractorSuite) TestDeleteExtractor_Execute_Returns_An_Error_If_The_Extractor_Cannot_Be_Deleted() {
	suite.client.ExpectedCalls = nil
	suite.client.On("GetAll", mock.Anything).Return(
		AListOfExtractors("charlie", "jo", "chris"), nil)

	expected := errors.New("Failed")
	suite.client.On("Delete", mock.Anything, mock.Anything).Return(expected)

	actual := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expected, actual)
}
