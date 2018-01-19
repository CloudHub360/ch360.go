package tests

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BasicFileSinkSuite struct {
	suite.Suite
	sut          *sinks.BasicFileSink
	fileSystem   afero.Fs
	filename     string
	filecontents string
}

func (suite *BasicFileSinkSuite) SetupTest() {
	suite.filename = generators.String("filename")
	suite.filecontents = generators.String("contents")
	suite.fileSystem = afero.NewMemMapFs()
	suite.sut = sinks.NewBasicFileSink(suite.fileSystem, suite.filename)
}

func TestBasicFileSinkRunner(t *testing.T) {
	suite.Run(t, new(BasicFileSinkSuite))
}

func (suite *BasicFileSinkSuite) TestWritesFiles() {
	// Act
	suite.sut.Open()
	fmt.Fprint(suite.sut, suite.filecontents)
	suite.sut.Close()

	// Assert
	contents, err := afero.ReadFile(suite.fileSystem, suite.filename)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.filecontents, fmt.Sprintf("%s", contents))
}
