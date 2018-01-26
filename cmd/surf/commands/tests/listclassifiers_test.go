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

type ListClassifierSuite struct {
	suite.Suite
	sut    *commands.ListClassifiers
	client *mocks.ClassifierGetter
	output *bytes.Buffer
}

func (suite *ListClassifierSuite) SetupTest() {
	suite.client = new(mocks.ClassifierGetter)
	suite.output = &bytes.Buffer{}

	suite.sut = commands.NewListClassifiers(suite.output, suite.client)
}

func TestListClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListClassifierSuite))
}

func (suite *ListClassifierSuite) TestGetAllClassifiers_Execute_Returns_The_Classifiers_When_There_Are_Some() {
	expectedClassifiers := AListOfClassifiers("charlie", "jo", "chris").(ch360.ClassifierList)
	suite.client.On("GetAll", mock.Anything).Return(expectedClassifiers, nil)

	classifiers, _ := suite.sut.Execute()

	suite.client.AssertCalled(suite.T(), "GetAll")
	assert.Equal(suite.T(), expectedClassifiers, classifiers)
}

func (suite *ListClassifierSuite) TestGetAllClassifiers_Execute_Returns_Empty_Classifiers_List_When_There_Are_None() {
	expectedClassifiers := make(ch360.ClassifierList, 0)
	suite.client.On("GetAll", mock.Anything).Return(expectedClassifiers, nil)

	classifiers, _ := suite.sut.Execute()

	suite.client.AssertCalled(suite.T(), "GetAll")
	assert.Equal(suite.T(), expectedClassifiers, classifiers)
}

func (suite *ListClassifierSuite) TestGetAllClassifiers_Execute_Returns_An_Error_If_The_Classifiers_Cannot_Be_Retrieved() {
	expectedErr := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expectedErr)

	_, actualErr := suite.sut.Execute()

	assert.Equal(suite.T(), expectedErr, actualErr)
}
