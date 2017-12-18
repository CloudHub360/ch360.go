package commands

import (
	"github.com/CloudHub360/ch360.go/config/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoginSuite struct {
	suite.Suite
	sut             *Login
	configDirectory *mocks.Writer
}

func (suite *LoginSuite) SetupTest() {
	suite.configDirectory = new(mocks.Writer)
	suite.configDirectory.On("Write", mock.Anything).Return(0, nil)

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

	suite.configDirectory.AssertCalled(suite.T(), "Write", mock.Anything)
}
