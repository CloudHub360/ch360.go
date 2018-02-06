package commands

import (
	"errors"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CredentialsResolverSuite struct {
	suite.Suite
	sut                *CredentialsResolver
	reader             *mocks.ConfigurationReader
	configClientId     string
	configClientSecret string
}

func (suite *CredentialsResolverSuite) SetupTest() {
	suite.configClientId = generators.String("config-clientid")
	suite.configClientSecret = generators.String("config-clientsecret")

	configuration := config.NewConfiguration(suite.configClientId, suite.configClientSecret)
	suite.reader = new(mocks.ConfigurationReader)
	suite.reader.On("ReadConfiguration").Return(configuration, nil)

	suite.sut = &CredentialsResolver{}
}

func TestCredentialsResolverSuiteRunner(t *testing.T) {
	suite.Run(t, new(CredentialsResolverSuite))
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Parameters_If_Both_Set() {
	clientIdParam := generators.String("clientid")
	clientSecretParam := generators.String("clientsecret")

	clientIdActual, clientSecretActual, err := suite.sut.Resolve(clientIdParam, clientSecretParam, suite.reader)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), clientIdParam, clientIdActual)
	assert.Equal(suite.T(), clientSecretParam, clientSecretActual)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_Id_Parameter_Set_But_Not_Secret() {
	clientIdParam := generators.String("clientid")
	clientSecretParam := ""
	expectedErr := errors.New("You must either specify both --client-id and --client-secret, or neither.")

	_, _, err := suite.sut.Resolve(clientIdParam, clientSecretParam, suite.reader)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_Secret_Parameter_Set_But_Not_Id() {
	clientIdParam := ""
	clientSecretParam := generators.String("clientsecret")
	expectedErr := errors.New("You must either specify both --client-id and --client-secret, or neither.")

	_, _, err := suite.sut.Resolve(clientIdParam, clientSecretParam, suite.reader)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Config_Values_If_Neither_Secret_Nor_Id_Set() {
	clientIdParam := ""
	clientSecretParam := ""

	clientIdActual, clientSecretActual, err := suite.sut.Resolve(clientIdParam, clientSecretParam, suite.reader)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.configClientId, clientIdActual)
	assert.Equal(suite.T(), suite.configClientSecret, clientSecretActual)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_Config_Values_Are_Needed_And_Id_Is_Empty() {
	configuration := config.NewConfiguration("", suite.configClientSecret)
	expectedErr := errors.New("Your configuration file does not contain valid credentials. Please run 'surf login' to connect to your account.")

	suite.reader.ExpectedCalls = nil
	suite.reader.On("ReadConfiguration").Return(configuration, nil)

	_, _, err := suite.sut.Resolve("", "", suite.reader)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_Config_Values_Are_Needed_And_Secret_Is_Empty() {
	configuration := config.NewConfiguration(suite.configClientId, "")
	expectedErr := errors.New("Your configuration file does not contain valid credentials. Please run 'surf login' to connect to your account.")

	suite.reader.ExpectedCalls = nil
	suite.reader.On("ReadConfiguration").Return(configuration, nil)

	_, _, err := suite.sut.Resolve("", "", suite.reader)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_Config_Values_Are_Needed_And_Contains_No_Credentials() {
	var credentials = make(config.ApiCredentialsList, 0)
	configuration := &config.Configuration{
		Credentials: credentials,
	}
	expectedErr := errors.New("Your configuration file does not contain any credentials. Please run 'surf login' to connect to your account.")

	suite.reader.ExpectedCalls = nil
	suite.reader.On("ReadConfiguration").Return(configuration, nil)

	_, _, err := suite.sut.Resolve("", "", suite.reader)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Error_If_ConfigurationReader_Returns_A_No_ConfigFile_Error() {
	expectedError := errors.New("Please run 'surf login' to connect to your account.")

	suite.reader.ExpectedCalls = nil
	suite.reader.On("ReadConfiguration").Return(nil, os.ErrNotExist)

	_, _, err := suite.sut.Resolve("", "", suite.reader)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *CredentialsResolverSuite) TestResolve_Returns_Wrapped_Error_If_ConfigurationReader_Returns_Another_Error() {
	errorText := "Corrupted file"
	configReadingError := errors.New(errorText)
	expectedError := errors.New("There was an error loading your configuration file. Please run 'surf login' to connect to your account. Error: " + errorText)

	suite.reader.ExpectedCalls = nil
	suite.reader.On("ReadConfiguration").Return(nil, configReadingError)

	_, _, err := suite.sut.Resolve("", "", suite.reader)
	assert.Equal(suite.T(), expectedError, err)
}
