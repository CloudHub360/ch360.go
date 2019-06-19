package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateExtractorSuite struct {
	suite.Suite
	output          *bytes.Buffer
	creator         *mocks.ExtractorCreator
	sut             *commands.CreateExtractor
	modulesTemplate *ch360.ExtractorTemplate
	extractorName   string
	ctx             context.Context
}

const modulesTemplateJson = `{"modules":[{"id":"waives.name"},{"id":"waives.date"}]}`

func (suite *CreateExtractorSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.creator = new(mocks.ExtractorCreator)

	suite.modulesTemplate = aModulesTemplate()
	suite.extractorName = generators.String("extractor-name")
	suite.sut = commands.NewCreateExtractor(suite.output,
		suite.creator,
		suite.extractorName,
		suite.modulesTemplate)

	suite.creator.On("CreateFromModules", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.ctx = context.Background()
}

func aModulesTemplate() *ch360.ExtractorTemplate {
	modulesTemplate, _ := ch360.NewModulesTemplateFromJson(bytes.NewBufferString(modulesTemplateJson))
	return modulesTemplate
}

func TestCreateExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateExtractorSuite))
}

func (suite *CreateExtractorSuite) ClearExpectedCalls() {
	suite.creator.ExpectedCalls = nil
}

func (suite *CreateExtractorSuite) TestCreateExtractor_Execute_Calls_Client_With_Correct_Args() {
	suite.sut.Execute(context.Background())

	suite.creator.AssertCalled(suite.T(), "CreateFromModules", suite.ctx, suite.extractorName, *suite.modulesTemplate)
}

func (suite *CreateExtractorSuite) TestCreateExtractor_Execute_Returns_Error_If_The_Extractor_Cannot_Be_Created() {
	expectedErr := errors.New("Error message")
	suite.ClearExpectedCalls()
	suite.creator.
		On("CreateFromModules", mock.Anything, mock.Anything, mock.Anything).
		Return(expectedErr)

	receivedErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}
