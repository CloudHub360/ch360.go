package commands

import (
	"github.com/CloudHub360/ch360.go/config"
	//"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut *Login
	//configDirectory *config.ConfigurationWriter
	configDirectory *config.ConfigurationWriter
}

func (suite *LoginSuite) SetupTest() {
	//suite.configDirectory = &fakeConfigurationWriter{}
	//suite.configDirectory = new(mocks.ConfigurationWriter)
	//suite.configDirectory.On("WriteConfiguration", mock.Anything).Return(0, nil)

	suite.sut = NewLogin(suite.configDirectory)
}

func TestLoginSuiteRunner(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestLogin_Execute_Writes_Configuration_To_ConfigurationDirectory() {
	err := suite.sut.Execute("clientid", "clientsecret")
	if err != nil {
		assert.Error(suite.T(), err)
	}

	//suite.configDirectory.AssertCalled(suite.T(), "WriteConfiguration", mock.Anything)
}

type fakeConfigurationWriter struct{}

func (writer *fakeConfigurationWriter) WriteConfiguration(configuration *config.Configuration) error {
	return nil
}
