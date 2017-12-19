package tests

import (
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	cmdmocks "github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut          *commands.Login
	configWriter *mocks.ConfigurationWriter
	secretReader *cmdmocks.SecretReader
	clientId     string
	clientSecret string
}

func (suite *LoginSuite) SetupTest() {
	suite.clientId = generators.String("clientid")
	suite.clientSecret = generators.String("clientsecret")

	suite.configWriter = new(mocks.ConfigurationWriter)
	suite.configWriter.On("WriteConfiguration", mock.Anything).Return(nil)

	suite.secretReader = new(cmdmocks.SecretReader)
	suite.secretReader.On("Read").Return(suite.clientSecret, nil)
	suite.sut = commands.NewLogin(suite.configWriter, suite.secretReader)
}

func TestLoginSuiteRunner(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_When_Id_And_Secret_Specified() {
	err := suite.sut.Execute(suite.clientId, suite.clientSecret)
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.assertConfigurationWrittenWithCredentials(suite.clientId, suite.clientSecret)
}

func (suite *LoginSuite) TestLogin_Execute_Prompts_For_Secret_When_Secret_Not_Specified() {
	err := suite.sut.Execute(suite.clientId, "")
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.secretReader.AssertCalled(suite.T(), "Read")
}

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_When_Secret_Is_Entered_At_Prompt() {
	// Simulate user entering secret at prompt by the mock SecretReader returning the secret
	err := suite.sut.Execute(suite.clientId, "")
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.assertConfigurationWrittenWithCredentials(suite.clientId, suite.clientSecret)
}

func (suite *LoginSuite) assertConfigurationWrittenWithCredentials(clientId string, clientSecret string) {
	suite.configWriter.AssertCalled(suite.T(), "WriteConfiguration", mock.Anything)

	call := suite.configWriter.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)
	configuration := call.Arguments[0].(*config.Configuration)
	assert.Equal(suite.T(), clientId, configuration.ConfigurationRoot.Credentials[0].Id)
	assert.Equal(suite.T(), clientSecret, configuration.ConfigurationRoot.Credentials[0].Secret)
}
