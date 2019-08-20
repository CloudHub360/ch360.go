package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ListExtractorSuite struct {
	suite.Suite
	sut    *commands.ListExtractorsCmd
	client *mocks.ExtractorGetter
	output *bytes.Buffer
	ctx    context.Context
}

func (suite *ListExtractorSuite) SetupTest() {
	suite.client = new(mocks.ExtractorGetter)
	suite.output = &bytes.Buffer{}

	suite.sut = &commands.ListExtractorsCmd{
		Client: suite.client,
	}
	suite.ctx = context.Background()
}

func TestListExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListExtractorSuite))
}

func (suite *ListExtractorSuite) TestGetAllExtractors_Execute_Calls_The_Client() {
	expectedExtractors := AListOfExtractors("charlie", "jo", "chris").(ch360.ExtractorList)
	suite.client.On("GetAll", mock.Anything).Return(expectedExtractors, nil)

	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
}

func (suite *ListExtractorSuite) TestGetAllExtractors_Execute_Returns_An_Error_If_The_Extractors_Cannot_Be_Retrieved() {
	expectedErr := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expectedErr)

	actualErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func AListOfExtractors(names ...string) interface{} {
	expected := make(ch360.ExtractorList, len(names))

	for index, name := range names {
		expected[index] = ch360.Extractor{name}
	}

	return expected
}
