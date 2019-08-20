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

type ListClassifierSuite struct {
	suite.Suite
	sut    *commands.ListClassifiersCmd
	client *mocks.ClassifierGetter
	output *bytes.Buffer
	ctx    context.Context
}

func (suite *ListClassifierSuite) SetupTest() {
	suite.client = new(mocks.ClassifierGetter)
	suite.output = &bytes.Buffer{}

	suite.sut = &commands.ListClassifiersCmd{
		Client: suite.client,
	}
}

func TestListClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListClassifierSuite))
}

func (suite *ListClassifierSuite) TestGetAllClassifiers_Execute_Calls_Client() {
	expectedClassifiers := AListOfClassifiers("charlie", "jo", "chris").(ch360.ClassifierList)
	suite.client.On("GetAll", mock.Anything).Return(expectedClassifiers, nil)

	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
}

func (suite *ListClassifierSuite) TestGetAllClassifiers_Execute_Returns_Error_From_Client() {
	expectedClassifiers := make(ch360.ClassifierList, 0)
	expectedErr := errors.New("simulated err")
	suite.client.On("GetAll", mock.Anything).Return(expectedClassifiers, expectedErr)

	receivedErr := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expectedErr, receivedErr)
}
