package tests

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/output/sinks"
	"github.com/waives/surf/test/generators"
	"testing"
)

type ExtensionSwappingFileSinkSuite struct {
	suite.Suite
	sut                         *sinks.ExtensionSwappingFileSink
	fileSystem                  afero.Fs
	inputFilename               string
	expectedDestinationFilename string
	filecontents                string
	newExtension                string
}

func (suite *ExtensionSwappingFileSinkSuite) SetupTest() {
	filename := generators.String("inputFilename")
	suite.inputFilename = "/var/folder/" + filename + ".tif"
	suite.expectedDestinationFilename = "/var/folder/" + filename + ".ext"
	suite.newExtension = ".ext"
	suite.filecontents = generators.String("contents")
	suite.fileSystem = afero.NewMemMapFs()
	suite.sut = sinks.NewExtensionSwappingFileSink(suite.fileSystem, suite.newExtension, suite.inputFilename)
}

func TestExtensionSwappingFileSinkRunner(t *testing.T) {
	suite.Run(t, new(ExtensionSwappingFileSinkSuite))
}

func (suite *ExtensionSwappingFileSinkSuite) TestWritesFile_With_Correct_Filename_And_Contents() {
	// Act
	suite.sut.Open()
	fmt.Fprint(suite.sut, suite.filecontents)
	suite.sut.Close()

	// Assert
	fileExists, _ := afero.Exists(suite.fileSystem, suite.expectedDestinationFilename)
	require.True(suite.T(), fileExists)

	contents, err := afero.ReadFile(suite.fileSystem, suite.expectedDestinationFilename)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.filecontents, fmt.Sprintf("%s", contents))
}
