package config

import (
	"encoding/json"
	mockconfig "github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigurationSuite struct {
	suite.Suite
	sut                 *Configuration
	mockConfigDirectory *mockconfig.FileWriter
	clientId            string
	clientSecret        string
}

func (suite *ConfigurationSuite) SetupTest() {
	suite.mockConfigDirectory = &mockconfig.FileWriter{}
	suite.mockConfigDirectory.On("WriteFile", mock.Anything, mock.Anything).Return(nil)

	suite.clientId = "clientid"
	suite.clientSecret = "clientsecret"
	suite.sut = NewConfiguration(suite.clientId, suite.clientSecret)
}

func TestConfigurationSuiteRunner(t *testing.T) {
	suite.Run(t, new(ConfigurationSuite))
}

func (suite *ConfigurationSuite) TestConfigurationNewConfiguration_Creates_A_Configuration_With_Specified_Credentials() {
	// Assert
	assert.Len(suite.T(), suite.sut.ConfigurationRoot.Credentials, 1)
	actualCredentials := suite.sut.ConfigurationRoot.Credentials[0]
	assert.Equal(suite.T(), suite.clientId, actualCredentials.Id)
	assert.Equal(suite.T(), suite.clientSecret, actualCredentials.Secret)
	assert.Equal(suite.T(), "default", actualCredentials.Key)
	assert.Equal(suite.T(), "default", actualCredentials.Url)
}

func (suite *ConfigurationSuite) TestConfigurationSaves_Writes_File_With_Correct_Name() {
	// Act
	suite.sut.Save(suite.mockConfigDirectory)

	// Assert
	suite.mockConfigDirectory.AssertNumberOfCalls(suite.T(), "WriteFile", 1)

	call := suite.mockConfigDirectory.Calls[0]
	assert.Len(suite.T(), call.Arguments, 2)
	assert.Equal(suite.T(), "config.json", call.Arguments[0])
}

func (suite *ConfigurationSuite) TestConfigurationSaves_Writes_File_With_Serialised_Configuration() {
	// Act
	suite.sut.Save(suite.mockConfigDirectory)

	// Assert
	suite.mockConfigDirectory.AssertNumberOfCalls(suite.T(), "WriteFile", 1)

	call := suite.mockConfigDirectory.Calls[0]
	assert.Len(suite.T(), call.Arguments, 2)
	configuration := suite.AssertIsValidSerialisedConfiguration(call.Arguments[1].([]byte))
	suite.AssertConfigurationIsPopulatedWithData(configuration)
}

func (suite *ConfigurationSuite) AssertIsValidSerialisedConfiguration(contents []byte) Configuration {
	var configuration Configuration
	err := json.Unmarshal(contents, &configuration)
	if err != nil {
		assert.Fail(suite.T(), "Bytes are not valid serialised Configuration")
	}
	return configuration
}

func (suite *ConfigurationSuite) AssertConfigurationIsPopulatedWithData(configuration Configuration) {
	assert.Equal(suite.T(), suite.clientId, configuration.ConfigurationRoot.Credentials[0].Id)
}
