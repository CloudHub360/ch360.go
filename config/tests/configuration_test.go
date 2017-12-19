package tests

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/test/generators"
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
	assert.Len(suite.T(), suite.sut.Credentials, 1)
	actualCredentials := suite.sut.Credentials[0]
	assert.Equal(suite.T(), suite.clientId, actualCredentials.Id)
	assert.Equal(suite.T(), suite.clientSecret, actualCredentials.Secret)
	assert.Equal(suite.T(), "default", actualCredentials.Key)
	assert.Equal(suite.T(), "default", actualCredentials.Url)
}

func (suite *ConfigurationSuite) TestConfigurationSerialise_Can_Be_Deserialised_To_Configuration_With_Same_Values() {
	bytes, err := suite.sut.Serialise()
	assert.Nil(suite.T(), err)

	configuration, err := config.DeserialiseConfiguration(bytes)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), len(suite.sut.Credentials), len(configuration.Credentials))
	credentials := suite.sut.Credentials[0]
	assert.Equal(suite.T(), credentials.Id, suite.clientId)
	assert.Equal(suite.T(), credentials.Secret, suite.clientSecret)
	assert.Equal(suite.T(), credentials.Key, "default")
	assert.Equal(suite.T(), credentials.Url, "default")
}

func (suite *ConfigurationSuite) TestConfigurationDeserialise_Returns_Error_If_Attempting_To_Deserialised_Invalid_Contents() {
	_, err := config.DeserialiseConfiguration(generators.Bytes())
	assert.NotNil(suite.T(), err)
}
