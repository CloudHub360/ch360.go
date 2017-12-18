package tests

import (
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	cmdmocks "github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut             *commands.Login
	configDirectory *mocks.ConfigurationWriter
	secretReader    *cmdmocks.SecretReader
	clientId        string
	clientSecret    string
}

func (suite *LoginSuite) SetupTest() {
	suite.configDirectory = new(mocks.ConfigurationWriter)
	suite.configDirectory.On("WriteConfiguration", mock.Anything).Return(nil)

	suite.secretReader = new(cmdmocks.SecretReader)
	suite.secretReader.On("Read").Return(suite.clientSecret, nil)
	suite.sut = commands.NewLogin(suite.configDirectory, suite.secretReader)

	suite.clientId = suite.clientId
	suite.clientSecret = suite.clientSecret
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

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_When_Secret_Not_Specified() {
	err := suite.sut.Execute(suite.clientId, "")
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.assertConfigurationWrittenWithCredentials(suite.clientId, suite.clientSecret)
}

func (suite *LoginSuite) assertConfigurationWrittenWithCredentials(clientId string, clientSecret string) {
	suite.configDirectory.AssertCalled(suite.T(), "WriteConfiguration", mock.Anything)

	call := suite.configDirectory.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)
	configuration := call.Arguments[0].(*config.Configuration)
	assert.Equal(suite.T(), clientId, configuration.ConfigurationRoot.Credentials[0].Id)
	assert.Equal(suite.T(), clientSecret, configuration.ConfigurationRoot.Credentials[0].Secret)
}
