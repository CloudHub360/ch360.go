package config

import (
	"fmt"
	fakes "github.com/CloudHub360/ch360.go/config/fakes"
	assertThat "github.com/CloudHub360/ch360.go/test/assertions"
	generate "github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
	"time"
)

type ConfigurationDirectorySuite struct {
	suite.Suite
	sut                   *ConfigurationDirectory
	homeDirectory         *fakes.FakeHomeDirectoryPathGetter
	fileContents          []byte
	expectedDirectoryPath string
	expectedFilename      string
	expectedFilePath      string
}

func (suite *ConfigurationDirectorySuite) SetupTest() {
	// Create unique "home directory" for this test
	suite.homeDirectory = &fakes.FakeHomeDirectoryPathGetter{
		Guid: fmt.Sprintf("%v", time.Now().UTC().UnixNano()),
	}
	suite.homeDirectory.Create()

	suite.sut = NewConfigurationDirectory(
		suite.homeDirectory)

	suite.expectedFilename = "config.yaml"
	suite.expectedDirectoryPath = filepath.Join(
		suite.homeDirectory.GetPath(),
		".ch360")
	suite.expectedFilePath = filepath.Join(
		suite.expectedDirectoryPath,
		suite.expectedFilename)
	suite.fileContents = generate.Bytes()

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
	_, err := suite.sut.Write(suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assertThat.DirectoryExists(suite.T(), suite.expectedDirectoryPath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Name() {
	// Act
	_, err := suite.sut.Write(suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assertThat.FileExists(suite.T(), suite.expectedFilePath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Content() {
	// Act
	_, err := suite.sut.Write(suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	assertThat.FileHasContents(suite.T(), suite.expectedFilePath, suite.fileContents)
}
