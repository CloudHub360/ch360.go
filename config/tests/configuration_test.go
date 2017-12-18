package tests

import (
	"github.com/CloudHub360/ch360.go/config"
	mockconfig "github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"testing"
)

type ConfigurationSuite struct {
	suite.Suite
	sut                 *config.Configuration
	mockConfigDirectory *mockconfig.Writer
	clientId            string
	clientSecret        string
}

func (suite *ConfigurationSuite) SetupTest() {
	suite.mockConfigDirectory = &mockconfig.Writer{}
	suite.mockConfigDirectory.On("Write", mock.Anything).Return(0, nil)

	suite.clientId = "clientid"
	suite.clientSecret = "clientsecret"
	suite.sut = config.NewConfiguration(suite.clientId, suite.clientSecret)
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

func (suite *ConfigurationSuite) TestConfigurationSaves_Writes_File_With_Serialised_Configuration() {
	// Act
	suite.sut.Save(suite.mockConfigDirectory)

	// Assert
	suite.mockConfigDirectory.AssertNumberOfCalls(suite.T(), "Write", 1)

	call := suite.mockConfigDirectory.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)
	configuration := suite.AssertIsValidSerialisedConfiguration(call.Arguments[0].([]byte))
	suite.AssertConfigurationIsPopulatedWithData(configuration)
}

func (suite *ConfigurationSuite) AssertIsValidSerialisedConfiguration(contents []byte) config.Configuration {
	var configuration config.Configuration
	err := yaml.Unmarshal(contents, &configuration)
	if err != nil {
		assert.Fail(suite.T(), "Bytes are not valid serialised Configuration")
	}
	return configuration
}

func (suite *ConfigurationSuite) AssertConfigurationIsPopulatedWithData(configuration config.Configuration) {
	assert.Equal(suite.T(), suite.clientId, configuration.ConfigurationRoot.Credentials[0].Id)
}
