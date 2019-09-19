package tests

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/cmd/surf/commands"
	"github.com/waives/surf/cmd/surf/commands/mocks"
	"testing"
)

type ListModuleSuite struct {
	suite.Suite
	sut    *commands.ListModulesCmd
	client *mocks.ModuleGetter
	ctx    context.Context
}

func (suite *ListModuleSuite) SetupTest() {
	suite.client = new(mocks.ModuleGetter)

	suite.sut = &commands.ListModulesCmd{
		Client: suite.client,
	}
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

func aListOfModules(ids ...string) interface{} {
	expected := make(ch360.ModuleList, len(ids))

	for index, id := range ids {
		expected[index] = ch360.Module{
			Name: id,
			ID:   id,
		}
	}

	return expected
}
