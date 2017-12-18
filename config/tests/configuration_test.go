package tests

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigurationSuite struct {
	suite.Suite
	sut          *config.Configuration
	clientId     string
	clientSecret string
}

func (suite *ConfigurationSuite) SetupTest() {
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

func (suite *ConfigurationSuite) TestConfigurationSerialise_Can_Be_Deserialised_To_Configuration_With_Same_Values() {
	bytes, err := suite.sut.Serialise()
	if err != nil {
		assert.Error(suite.T(), err)
	}

	configuration, err := config.DeserialiseConfiguration(bytes)
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assert.Equal(suite.T(), len(suite.sut.ConfigurationRoot.Credentials), len(configuration.ConfigurationRoot.Credentials))
	credentials := suite.sut.ConfigurationRoot.Credentials[0]
	assert.Equal(suite.T(), credentials.Id, suite.clientId)
	assert.Equal(suite.T(), credentials.Secret, suite.clientSecret)
	assert.Equal(suite.T(), credentials.Key, "default")
	assert.Equal(suite.T(), credentials.Url, "default")
}
