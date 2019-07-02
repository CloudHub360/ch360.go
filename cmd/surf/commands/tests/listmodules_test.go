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

type ListModuleSuite struct {
	suite.Suite
	sut    *commands.ListModules
	client *mocks.ModuleGetter
	output *bytes.Buffer
	ctx    context.Context
}

func (suite *ListModuleSuite) SetupTest() {
	suite.client = new(mocks.ModuleGetter)
	suite.output = &bytes.Buffer{}

	suite.sut = commands.NewListModules(suite.client, suite.output)
	suite.ctx = context.Background()
}

func TestListModuleSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListModuleSuite))
}

func (suite *ListModuleSuite) TestGetAllModules_Execute_Calls_The_Client() {
	expectedModules := aListOfModules("charlie", "jo", "chris").(ch360.ModuleList)
	suite.client.On("GetAll", mock.Anything).Return(expectedModules, nil)

	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
}

func (suite *ListModuleSuite) TestGetAllModules_Execute_Returns_An_Error_If_The_Modules_Cannot_Be_Retrieved() {
	expectedErr := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expectedErr)

	actualErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func aListOfModules(names ...string) interface{} {
	expected := make(ch360.ModuleList, len(names))

	for index, name := range names {
		expected[index] = ch360.Module{
			Name: name,
			ID:   name,
		}
	}

	return expected
}
