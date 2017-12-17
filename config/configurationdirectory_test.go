package config

import (
	"fmt"
	fakes "github.com/CloudHub360/ch360.go/config/fakes"
	assertThat "github.com/CloudHub360/ch360.go/test/assertions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type ConfigurationDirectorySuite struct {
	suite.Suite
	sut                   *ConfigurationDirectory
	fileSystem            *FileSystem
	homeDirectory         *fakes.FakeHomeDirectoryPathGetter
	fileContents          []byte
	expectedDirectoryPath string
	filename              string
	expectedFilePath      string
}

func (suite *ConfigurationDirectorySuite) SetupTest() {
	// Create unique "home directory" for this test
	suite.homeDirectory = &fakes.FakeHomeDirectoryPathGetter{
		Guid: fmt.Sprintf("%v", time.Now().UTC().UnixNano()),
	}
	suite.homeDirectory.Create()

	suite.sut = NewConfigurationDirectory(
		suite.homeDirectory,
		suite.fileSystem)

	suite.filename = "assertThat-config-file.json"
	suite.expectedDirectoryPath = suite.fileSystem.JoinPath(
		suite.homeDirectory.GetPath(),
		".ch360")
	suite.expectedFilePath = suite.fileSystem.JoinPath(
		suite.expectedDirectoryPath,
		suite.filename)
	suite.fileContents = generateBytes()

	assertThat.DirectoryDoesNotExist(suite.T(), suite.expectedDirectoryPath)
	assertThat.FileDoesNotExist(suite.T(), suite.expectedFilePath)
}

func (suite *ConfigurationDirectorySuite) TearDownTest() {
	suite.homeDirectory.Destroy()
}

func TestConfigurationDirectorySuiteRunner(t *testing.T) {
	suite.Run(t, new(ConfigurationDirectorySuite))
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_Config_Directory_If_None_Exists() {
	// Act
	suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	assertThat.DirectoryExists(suite.T(), suite.expectedDirectoryPath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Name() {
	// Act
	err := suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assertThat.FileExists(suite.T(), suite.expectedFilePath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Content() {
	// Act
	err := suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assertThat.FileHasContents(suite.T(), suite.expectedFilePath, suite.fileContents)
}

func generateBytes() []byte {
	token := make([]byte, 100)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(token)
	return token
}
