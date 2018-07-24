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

type DeleteClassifierSuite struct {
	suite.Suite
	sut            *commands.DeleteClassifier
	output         *bytes.Buffer
	client         *mocks.ClassifierDeleterGetter
	classifierName string
	ctx            context.Context
}

func (suite *DeleteClassifierSuite) SetupTest() {
	suite.client = new(mocks.ClassifierDeleterGetter)
	suite.client.On("GetAll", mock.Anything).Return(
		AListOfClassifiers("charlie", "jo", "chris"), nil)

	suite.client.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.output = &bytes.Buffer{}

	suite.classifierName = "charlie"
	suite.sut = commands.NewDeleteClassifier(suite.classifierName, suite.output, suite.client)
	suite.ctx = context.Background()
}

func TestDeleteClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(DeleteClassifierSuite))
}

func (suite *DeleteClassifierSuite) TestDeleteClassifier_Execute_Deletes_The_Named_Classifier_When_It_Exists() {
	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
	suite.client.AssertCalled(suite.T(), "Delete", suite.ctx, suite.classifierName)
}

func (suite *DeleteClassifierSuite) TestDeleteClassifier_Execute_Does_Not_Delete_The_Named_Classifier_When_It_Does_Not_Exist() {
	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
	suite.client.AssertNotCalled(suite.T(), "Delete", suite.ctx)
}

func (suite *DeleteClassifierSuite) TestDeleteClassifier_Execute_Returns_An_Error_If_The_Classifiers_Cannot_Be_Retrieved() {
	suite.client.ExpectedCalls = nil
	expected := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything, mock.Anything).Return(nil, expected)

	actual := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expected, actual)
	suite.client.AssertNotCalled(suite.T(), "Delete")
}

func (suite *DeleteClassifierSuite) TestDeleteClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Deleted() {
	suite.client.ExpectedCalls = nil
	suite.client.On("GetAll", mock.Anything, mock.Anything).Return(
		AListOfClassifiers("charlie", "jo", "chris"), nil)

	expected := errors.New("Failed")
	suite.client.On("Delete", mock.Anything, mock.Anything).Return(expected)

	actual := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expected, actual)
}

func AListOfClassifiers(names ...string) interface{} {
	expected := make(ch360.ClassifierList, len(names))

	for index, name := range names {
		expected[index] = ch360.Classifier{name}
	}

	return expected
}
