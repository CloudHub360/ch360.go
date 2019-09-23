package tests

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/cmd/surf/commands"
	"github.com/waives/surf/cmd/surf/commands/mocks"
	"strings"
	"testing"
)

type CreateExtractorTemplateSuite struct {
	suite.Suite
	sut       *commands.CreateExtractorTemplateCmd
	client    *mocks.ModuleGetter
	moduleIds []string
	output    *bytes.Buffer
	ctx       context.Context
}

func (suite *CreateExtractorTemplateSuite) SetupTest() {
	suite.moduleIds = []string{"moduleA", "moduleB", "moduleC"}
	suite.client = new(mocks.ModuleGetter)
	suite.client.On("GetAll", mock.Anything).
		Return(aListOfModules(suite.moduleIds...), nil)

	suite.output = &bytes.Buffer{}

	suite.sut = &commands.CreateExtractorTemplateCmd{
		Client:    suite.client,
		ModuleIds: suite.moduleIds,
		Output:    suite.output,
	}
	suite.ctx = context.Background()
}

func TestCreateExtractorTemplateSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateExtractorTemplateSuite))
}

func (suite *CreateExtractorTemplateSuite) Test_CreateExtractorTemplate_Returns_Error_If_No_Modules_Are_Specified() {
	fixtures := []struct {
		moduleIds []string
	}{
		{
			moduleIds: []string{},
		}, {
			moduleIds: nil,
		},
	}

	for _, fixture := range fixtures {
		// Arrange
		sut := &commands.CreateExtractorTemplateCmd{
			Client:    suite.client,
			ModuleIds: fixture.moduleIds,
			Output:    suite.output,
		}

		// Act
		err := sut.Execute(suite.ctx)

		// Assert
		assert.Error(suite.T(), err)
	}
}

func (suite *CreateExtractorTemplateSuite) Test_CreateExtractorTemplate_Returns_Error_If_Specified_Modules_Do_Not_Exist() {

	// Arrange
	moduleIds := []string{"missingModuleA", "missingModuleB", "missingModuleC"}
	sut := &commands.CreateExtractorTemplateCmd{
		Client:    suite.client,
		ModuleIds: moduleIds,
		Output:    suite.output,
	}

	// Act
	err := sut.Execute(suite.ctx)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CreateExtractorTemplateSuite) Test_CreateExtractorTemplate_Writes_A_Valid_Template_To_Supplied_Writer() {
	// Act
	_ = suite.sut.Execute(suite.ctx)
	template, err := ch360.NewModulesTemplateFromJson(suite.output)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(template.Modules))
}

func (suite *CreateExtractorTemplateSuite) Test_CreateExtractorTemplate_Is_Case_Insensitive() {
	// Arrange
	moduleIds := suite.moduleIds
	for i, moduleID := range moduleIds {
		moduleIds[i] = strings.ToUpper(moduleID)
	}
	sut := &commands.CreateExtractorTemplateCmd{
		Client:    suite.client,
		ModuleIds: moduleIds,
		Output:    suite.output,
	}

	// Act
	err := sut.Execute(suite.ctx)
	assert.NoError(suite.T(), err)
	template, err := ch360.NewModulesTemplateFromJson(suite.output)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(template.Modules))
}
