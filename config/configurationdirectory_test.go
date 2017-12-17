package config

import (
	"fmt"
	fakes "github.com/CloudHub360/ch360.go/config/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"os"
	"testing"
	"time"
)

type ConfigurationDirectorySuite struct {
	suite.Suite
	sut                   *ConfigurationDirectory
	fileSystem            *FileSystem
	homeDirectory         string
	fileContents          []byte
	expectedDirectoryPath string
	filename              string
	expectedFilePath      string
}

func (suite *ConfigurationDirectorySuite) SetupTest() {
	// Create a unique "home directory" for this test
	fakeHomeDirectoryGetter := &fakes.FakeHomeDirectoryPathGetter{
		Guid: fmt.Sprintf("%v", time.Now().UTC().UnixNano()),
	}
	suite.homeDirectory = fakeHomeDirectoryGetter.GetPath()

	suite.fileSystem = &FileSystem{}
	suite.fileSystem.CreateDirectoryIfNotExists(suite.homeDirectory)

	suite.sut = NewConfigurationDirectory(
		fakeHomeDirectoryGetter,
		&FileSystem{})

	suite.filename = "a-config-file.json"
	suite.expectedDirectoryPath = suite.fileSystem.JoinPath(
		suite.homeDirectory,
		".ch360")
	suite.expectedFilePath = suite.fileSystem.JoinPath(
		suite.expectedDirectoryPath,
		suite.filename)
	suite.fileContents = generateBytes()

	suite.assertDirectoryDoesNotExist(suite.fileSystem, suite.expectedDirectoryPath)
	suite.assertFileDoesNotExist(suite.fileSystem, suite.expectedFilePath)

}

func (suite *ConfigurationDirectorySuite) TearDownTest() {
	os.RemoveAll(suite.homeDirectory)
}

func TestConfigurationDirectorySuiteRunner(t *testing.T) {
	suite.Run(t, new(ConfigurationDirectorySuite))
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_Config_Directory_If_None_Exists() {
	// Act
	suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	suite.assertDirectoryExists(suite.fileSystem, suite.expectedDirectoryPath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Name() {
	// Act
	err := suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.assertFileExists(suite.fileSystem, suite.expectedFilePath)
}

func (suite *ConfigurationDirectorySuite) TestConfigurationDirectoryWriteFile_Creates_File_With_Correct_Content() {
	// Act
	err := suite.sut.WriteFile(suite.filename, suite.fileContents)

	// Assert
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.assertFileHasContents(suite.fileSystem, suite.expectedFilePath, suite.fileContents)
}

func (suite *ConfigurationDirectorySuite) assertFileExists(fs *FileSystem, name string) {
	//TODO Change to FileExists
	exists, _ := fs.DirectoryExists(name)
	if !exists {
		assert.Fail(suite.T(), fmt.Sprintf("File %s does not exist", name))
	}
}

func (suite *ConfigurationDirectorySuite) assertFileDoesNotExist(fs *FileSystem, name string) {
	//TODO Change to FileExists
	exists, _ := fs.DirectoryExists(name)
	if exists {
		assert.Fail(suite.T(), fmt.Sprintf("File %s exists when it should not", name))
	}
}

func (suite *ConfigurationDirectorySuite) assertDirectoryExists(fs *FileSystem, name string) {
	exists, _ := fs.DirectoryExists(name)
	if !exists {
		assert.Fail(suite.T(), fmt.Sprintf("Directory %s does not exist", name))
	}
}

func (suite *ConfigurationDirectorySuite) assertDirectoryDoesNotExist(fs *FileSystem, name string) {
	exists, _ := fs.DirectoryExists(name)
	if exists {
		assert.Fail(suite.T(), fmt.Sprintf("Directory %s exists when it should not", name))
	}
}

func (suite *ConfigurationDirectorySuite) assertFileHasContents(fs *FileSystem, name string, contents []byte) {
	suite.assertFileExists(fs, name)
	contents, err := fs.ReadFile(name)
	if err != nil {
		assert.Error(suite.T(), err)
	}
	assert.Equal(suite.T(), suite.fileContents, contents)
}

func generateBytes() []byte {
	token := make([]byte, 100)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(token)
	return token
}
