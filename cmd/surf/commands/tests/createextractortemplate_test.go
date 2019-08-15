package tests

import (
	"bytes"
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CreateExtractorTemplateSuite struct {
	suite.Suite
	sut       *commands.CreateExtractorTemplate
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

	suite.sut = commands.NewCreateExtractorTemplate(suite.moduleIds, suite.client, suite.output)
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
		sut := commands.NewCreateExtractorTemplate(fixture.moduleIds, suite.client, suite.output)

		// Act
		err := sut.Execute(suite.ctx)

		// Assert
		assert.Error(suite.T(), err)
	}
}

func (suite *CreateExtractorTemplateSuite) Test_CreateExtractorTemplate_Returns_Error_If_Specified_Modules_Do_Not_Exist() {

	// Arrange
	moduleIds := []string{"missingModuleA", "missingModuleB", "missingModuleC"}
	sut := commands.NewCreateExtractorTemplate(moduleIds, suite.client, suite.output)

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
	sut := commands.NewCreateExtractorTemplate(moduleIds, suite.client, suite.output)

	// Act
	_ = sut.Execute(suite.ctx)
	template, err := ch360.NewModulesTemplateFromJson(suite.output)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(template.Modules))
}
