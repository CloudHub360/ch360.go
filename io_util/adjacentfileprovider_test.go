package io_util_test

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const extension = ".test"

type AdjacentFileProviderSuite struct {
	suite.Suite
	sut         *io_util.AdjacentFileProvider
	tmpFile     *os.File
	createdFile *os.File
}

func (suite *AdjacentFileProviderSuite) TearDownTest() {
	os.Remove(suite.tmpFile.Name())
	os.Remove(suite.createdFile.Name())
}
func (suite *AdjacentFileProviderSuite) SetupTest() {
	suite.sut = &io_util.AdjacentFileProvider{Extension: extension}
	suite.tmpFile, _ = ioutil.TempFile("", "AdjacentFileProviderSuite")
	suite.createdFile, _ = suite.sut.Provide(suite.tmpFile.Name())
}

func TestAdjacentFileProviderSuiteRunner(t *testing.T) {
	suite.Run(t, new(AdjacentFileProviderSuite))
}

func (suite *AdjacentFileProviderSuite) Test_AdjacentFileProvider_Creates_File_In_Same_Dir_As_Source() {
	assert.True(suite.T(), fileExists(suite.createdFile.Name()))
	assert.Equal(suite.T(), filepath.Dir(suite.tmpFile.Name()), filepath.Dir(suite.createdFile.Name()))
}

func (suite *AdjacentFileProviderSuite) Test_AdjacentFileProvider_Creates_File_With_Right_Extension() {
	assert.Equal(suite.T(), extension, filepath.Ext(suite.createdFile.Name()))
}

func (suite *AdjacentFileProviderSuite) Test_AdjacentFileProvider_Creates_File_With_Identical_Path_Except_Extension() {
	tmpPathWithoutExt := strings.TrimSuffix(suite.tmpFile.Name(), filepath.Ext(suite.tmpFile.Name()))
	createdPathWithoutExt := strings.TrimSuffix(suite.createdFile.Name(), filepath.Ext(suite.createdFile.Name()))

	assert.Equal(suite.T(), tmpPathWithoutExt, createdPathWithoutExt)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
