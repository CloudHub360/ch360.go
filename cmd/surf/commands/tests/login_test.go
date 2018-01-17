package tests

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/auth"
	authmocks "github.com/CloudHub360/ch360.go/auth/mocks"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut            *commands.Login
	configWriter   *mocks.ConfigurationWriter
	tokenRetriever *authmocks.TokenRetriever
	clientId       string
	clientSecret   string
	output         *bytes.Buffer
}

func (suite *LoginSuite) SetupTest() {
	suite.clientId = generators.String("clientid")
	suite.clientSecret = generators.String("clientsecret")

	suite.configWriter = new(mocks.ConfigurationWriter)
	suite.configWriter.On("WriteConfiguration", mock.Anything).Return(nil)

	suite.tokenRetriever = new(authmocks.TokenRetriever)
	suite.tokenRetriever.On("RetrieveToken").Return(&auth.AccessToken{}, nil)

	suite.output = &bytes.Buffer{}
	suite.sut = commands.NewLogin(suite.output, suite.configWriter, suite.tokenRetriever)
}

func TestLoginSuiteRunner(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_When_Id_And_Secret_Specified() {
	err := suite.sut.Execute(suite.clientId, suite.clientSecret)
	assert.Nil(suite.T(), err)

	suite.assertConfigurationWrittenWithCredentials(suite.clientId, suite.clientSecret)
}

func (suite *LoginSuite) TestLogin_Execute_Requests_Auth_Token() {
	err := suite.sut.Execute(suite.clientId, suite.clientSecret)
	assert.Nil(suite.T(), err)

	suite.tokenRetriever.AssertCalled(suite.T(), "RetrieveToken")
}

func (suite *LoginSuite) TestLogin_Execute_Returns_Err_From_Token_Retriever() {
	// Arrange
	suite.tokenRetriever.ExpectedCalls = nil

	expectedErr := errors.New("An error")
	suite.tokenRetriever.On("RetrieveToken").Return(nil, expectedErr)

	// Act
	err := suite.sut.Execute(suite.clientId, suite.clientSecret)

	// Assert
	suite.tokenRetriever.AssertCalled(suite.T(), "RetrieveToken")
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
