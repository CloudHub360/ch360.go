package tests

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/auth"
	authmocks "github.com/waives/surf/auth/mocks"
	"github.com/waives/surf/cmd/surf/commands"
	"github.com/waives/surf/config"
	"github.com/waives/surf/config/mocks"
	"github.com/waives/surf/test/generators"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut            *commands.LoginCmd
	configWriter   *mocks.ConfigurationWriter
	tokenRetriever *authmocks.TokenRetriever
	flags          *config.GlobalFlags
	output         *bytes.Buffer
	clientSecret   string
	clientId       string
}

func (suite *LoginSuite) SetupTest() {

	suite.clientId = generators.String("clientid")
	suite.clientSecret = generators.String("clientsecret")

	suite.configWriter = new(mocks.ConfigurationWriter)
	suite.configWriter.On("WriteConfiguration", mock.Anything).Return(nil)

	suite.tokenRetriever = new(authmocks.TokenRetriever)
	suite.tokenRetriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(&auth.AccessToken{}, nil)

	suite.flags = &config.GlobalFlags{
		ClientId:     suite.clientId,
		ClientSecret: suite.clientSecret,
	}

	suite.output = &bytes.Buffer{}
	suite.sut = &commands.LoginCmd{
		TokenRetriever:      suite.tokenRetriever,
		ConfigurationWriter: suite.configWriter,
	}
}

func TestLoginSuiteRunner(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_When_Id_And_Secret_Specified() {
	err := suite.sut.Execute(context.Background(), suite.flags)
	assert.Nil(suite.T(), err)

	suite.assertConfigurationWrittenWithCredentials(suite.flags.ClientId, suite.clientSecret)
}

func (suite *LoginSuite) TestLogin_Execute_Requests_Auth_Token() {
	err := suite.sut.Execute(context.Background(), suite.flags)
	assert.Nil(suite.T(), err)

	suite.tokenRetriever.AssertCalled(suite.T(), "RetrieveToken", suite.clientId, suite.clientSecret)
}

func (suite *LoginSuite) TestLogin_Execute_Returns_Err_From_Token_Retriever() {
	// Arrange
	suite.tokenRetriever.ExpectedCalls = nil

	expectedErr := errors.New("An error")
	suite.tokenRetriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Act
	err := suite.sut.Execute(context.Background(), suite.flags)

	// Assert
	suite.tokenRetriever.AssertCalled(suite.T(), "RetrieveToken", suite.clientId, suite.clientSecret)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *LoginSuite) assertConfigurationWrittenWithCredentials(clientId string, clientSecret string) {
	suite.configWriter.AssertCalled(suite.T(), "WriteConfiguration", mock.Anything)

	call := suite.configWriter.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)
	configuration := call.Arguments[0].(*config.Configuration)
	assert.Equal(suite.T(), clientId, configuration.Credentials[0].Id)
	assert.Equal(suite.T(), clientSecret, configuration.Credentials[0].Secret)
}
