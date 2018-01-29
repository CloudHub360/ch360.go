package tests

import (
	"bytes"
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
	sut    *commands.ListExtractors
	client *mocks.ExtractorGetter
	output *bytes.Buffer
}

func (suite *ListExtractorSuite) SetupTest() {
	suite.client = new(mocks.ExtractorGetter)
	suite.output = &bytes.Buffer{}

	suite.sut = commands.NewListExtractors(suite.output, suite.client)
}

func TestListExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListExtractorSuite))
}

func (suite *ListExtractorSuite) TestGetAllExtractors_Execute_Returns_The_Extractors_When_There_Are_Some() {
	expectedExtractors := AListOfExtractors("charlie", "jo", "chris").(ch360.ExtractorList)
	suite.client.On("GetAll", mock.Anything).Return(expectedExtractors, nil)

	extractors, _ := suite.sut.Execute()

	suite.client.AssertCalled(suite.T(), "GetAll")
	assert.Equal(suite.T(), expectedExtractors, extractors)
}

func (suite *ListExtractorSuite) TestGetAllExtractors_Execute_Returns_Empty_Extractors_List_When_There_Are_None() {
	expectedExtractors := make(ch360.ExtractorList, 0)
	suite.client.On("GetAll", mock.Anything).Return(expectedExtractors, nil)

	extractors, _ := suite.sut.Execute()

	suite.client.AssertCalled(suite.T(), "GetAll")
	assert.Equal(suite.T(), expectedExtractors, extractors)
}

func (suite *ListExtractorSuite) TestGetAllExtractors_Execute_Returns_An_Error_If_The_Extractors_Cannot_Be_Retrieved() {
	expectedErr := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expectedErr)

	_, actualErr := suite.sut.Execute()

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func AListOfExtractors(names ...string) interface{} {
	expected := make(ch360.ExtractorList, len(names))

	for index, name := range names {
		expected[index] = ch360.Extractor{name}
	}

	return expected
}