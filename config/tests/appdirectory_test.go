package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/config"
	assertThat "github.com/waives/surf/test/assertions"
	"github.com/waives/surf/test/generators"
	"path/filepath"
	"runtime"
	"testing"
)

type AppDirectorySuite struct {
	suite.Suite
	sut                    *config.AppDirectory
	config                 *config.Configuration
	homeDirectory          *TemporaryDirectory
	expectedConfigDir      string
	expectedConfigFilePath string
	clientId               string
	clientSecret           string
}

func (suite *AppDirectorySuite) SetupTest() {
	// Create unique "home directory" for this test
	suite.homeDirectory = NewTemporaryDirectory()
	suite.homeDirectory.Create()

	suite.sut = config.NewAppDirectoryInDir(suite.homeDirectory.Path())
	suite.clientId = generators.String("clientid")
	suite.clientSecret = generators.String("clientsecret")
	suite.config = config.NewConfiguration(suite.clientId, suite.clientSecret)

	suite.expectedConfigDir = filepath.Join(
		suite.homeDirectory.Path(),
		".surf")
	suite.expectedConfigFilePath = filepath.Join(
		suite.homeDirectory.Path(),
		".surf",
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

func (suite *AppDirectorySuite) TestAppDirectoryWriteConfiguration_Creates_App_Directory_If_None_Exists() {
	err := suite.sut.WriteConfiguration(suite.config)

	assert.Nil(suite.T(), err)
	assertThat.DirectoryExists(suite.T(), suite.expectedConfigDir)

	//Permissions are not set correctly on Windows, only linux (on Windows they are always 777)
	if runtime.GOOS != "windows" {
		assertThat.DirectoryOrFileHasPermissions(suite.T(), suite.expectedConfigDir, config.DirRWPermissions)
	}
}

func (suite *AppDirectorySuite) TestAppDirectoryWriteConfiguration_Creates_File_With_Correct_Name() {

	err := suite.sut.WriteConfiguration(suite.config)
	assert.Nil(suite.T(), err)
	assertThat.FileExists(suite.T(), suite.expectedConfigFilePath)

	//Permissions are not set correctly on Windows, only linux (on Windows they are always 777)
	if runtime.GOOS != "windows" {
		assertThat.DirectoryOrFileHasPermissions(suite.T(), suite.expectedConfigFilePath, config.FileRWPermissions)
	}
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
