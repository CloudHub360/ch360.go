package tests

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/config/fakes"
	assertThat "github.com/CloudHub360/ch360.go/test/assertions"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type AppDirectorySuite struct {
	suite.Suite
	sut                    *config.AppDirectory
	config                 *config.Configuration
	homeDirectory          *fakes.FakeHomeDirectoryPathGetter
	expectedConfigDir      string
	expectedConfigFilePath string
	clientId               string
	clientSecret           string
}

func (suite *AppDirectorySuite) SetupTest() {
	// Create unique "home directory" for this test
	suite.homeDirectory = fakes.NewFakeHomeDirectoryPathGetter()
	suite.homeDirectory.Create()

	suite.sut = config.NewAppDirectory(suite.homeDirectory.Path())
	suite.clientId = generators.String("clientid")
	suite.clientSecret = generators.String("clientsecret")
	suite.config = config.NewConfiguration(suite.clientId, suite.clientSecret)

	suite.expectedConfigDir = filepath.Join(
		suite.homeDirectory.Path(),
		".ch360")
	suite.expectedConfigFilePath = filepath.Join(
		suite.homeDirectory.Path(),
		".ch360",
		"config.yaml")

	assertThat.DirectoryDoesNotExist(suite.T(), suite.expectedConfigDir)
	assertThat.FileDoesNotExist(suite.T(), suite.expectedConfigFilePath)
}

func (suite *AppDirectorySuite) TearDownTest() {
	suite.homeDirectory.Destroy()
}

func TestAppDirectorySuiteRunner(t *testing.T) {
	suite.Run(t, new(AppDirectorySuite))
}

func (suite *AppDirectorySuite) TestAppDirectoryWriteConfiguration_Creates_Config_Directory_If_None_Exists() {
	err := suite.sut.WriteConfiguration(suite.config)

	assert.Nil(suite.T(), err)
	assertThat.DirectoryExists(suite.T(), suite.expectedConfigDir)
}

func (suite *AppDirectorySuite) TestAppDirectoryWriteConfiguration_Creates_File_With_Correct_Name() {

	err := suite.sut.WriteConfiguration(suite.config)
	assert.Nil(suite.T(), err)
	assertThat.FileExists(suite.T(), suite.expectedConfigFilePath)
}

func (suite *AppDirectorySuite) TestAppDirectoryWriteConfiguration_Creates_File_With_Correct_Content() {
	err := suite.sut.WriteConfiguration(suite.config)

	assert.Nil(suite.T(), err)
	reloadedConfig, err := suite.sut.ReadConfiguration()
	assertConfigurationHasCredentials(suite.T(), reloadedConfig, suite.clientId, suite.clientSecret)
}

func assertConfigurationHasCredentials(t *testing.T, configuration *config.Configuration, clientId string, clientSecret string) {
	assert.Equal(t, clientId, configuration.Credentials[0].Id)
	assert.Equal(t, clientSecret, configuration.Credentials[0].Secret)
}
